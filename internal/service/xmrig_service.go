package service

import (
	"bufio"
	"context"
	"crypto/tls"
	"embed"
	"encoding/json"
	"fmt"
	"go-wails/internal/models"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed xmrig-embedded/*
var embeddedXMRig embed.FS

// XMRigService XMRig服务
type XMRigService struct {
	cmd           *exec.Cmd
	isRunning     bool
	mutex         sync.RWMutex
	ctx           context.Context
	configSvc     *ConfigService
	startTime     time.Time
	logBuffer     []string
	maxLogLines   int
	poolConnected bool
}

// NewXMRigService 创建XMRig服务
func NewXMRigService() *XMRigService {
	return &XMRigService{
		configSvc:   NewConfigService(),
		maxLogLines: 500,
		logBuffer:   make([]string, 0, 500),
	}
}

// SetContext 设置上下文
func (s *XMRigService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *XMRigService) runSilent(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	_ = cmd.Run()
}

// getXMRigExecutable 获取XMRig可执行文件路径
func (s *XMRigService) getXMRigExecutable() (string, error) {
	arch := runtime.GOARCH
	var embeddedPath string

	if arch == "arm64" {
		embeddedPath = "xmrig-embedded/xmrig-windows-arm64/xmrig.exe"
	} else {
		embeddedPath = "xmrig-embedded/xmrig-windows-amd64/xmrig.exe"
	}

	// 创建临时目录
	tempDir := filepath.Join(os.TempDir(), "xmrig-runtime")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}

	// 提取 xmrig.exe
	xmrigPath := filepath.Join(tempDir, "xmrig.exe")
	if err := s.extractEmbeddedFile(embeddedPath, xmrigPath); err != nil {
		return "", err
	}

	// 提取 WinRing0x64.sys (必需的驱动文件)
	var driverPath string
	if arch == "arm64" {
		driverPath = "xmrig-embedded/xmrig-windows-arm64/WinRing0x64.sys"
	} else {
		driverPath = "xmrig-embedded/xmrig-windows-amd64/WinRing0x64.sys"
	}
	driverDest := filepath.Join(tempDir, "WinRing0x64.sys")
	if err := s.extractEmbeddedFile(driverPath, driverDest); err != nil {
		// 驱动文件不是必需的，忽略错误
		fmt.Printf("警告: 无法提取驱动文件: %v\n", err)
	}

	return xmrigPath, nil
}

// extractEmbeddedFile 从嵌入的文件系统中提取文件
func (s *XMRigService) extractEmbeddedFile(embeddedPath, destPath string) error {
	// 如果文件已存在且大小正确，跳过提取
	if info, err := os.Stat(destPath); err == nil {
		// 检查嵌入文件的大小
		embeddedData, err := embeddedXMRig.ReadFile(embeddedPath)
		if err == nil && info.Size() == int64(len(embeddedData)) {
			return nil // 文件已存在且大小正确
		}
	}

	// 读取嵌入的文件
	data, err := embeddedXMRig.ReadFile(embeddedPath)
	if err != nil {
		return fmt.Errorf("读取嵌入文件失败 %s: %w", embeddedPath, err)
	}

	// 写入到目标路径
	if err := os.WriteFile(destPath, data, 0755); err != nil {
		return fmt.Errorf("写入文件失败 %s: %w", destPath, err)
	}

	return nil
}

// Start 启动挖矿
func (s *XMRigService) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isRunning {
		return fmt.Errorf("挖矿程序已在运行中")
	}

	exePath, err := s.getXMRigExecutable()
	if err != nil {
		return err
	}

	cfg, err := s.configSvc.LoadConfig()
	if err == nil && len(cfg.Pools) > 0 {
		// 选择第一个可用矿池，并重排到首位
		chosen := -1
		for i, p := range cfg.Pools {
			if !p.Enabled {
				continue
			}
			if s.isPoolReachable(p.URL) {
				chosen = i
				break
			}
		}
		if chosen == -1 {
			return fmt.Errorf("没有可用的矿池，请检查配置是否正确")
		}
		if chosen != 0 {
			pools := make([]models.PoolConfig, 0, len(cfg.Pools))
			pools = append(pools, cfg.Pools[chosen])
			for i, p := range cfg.Pools {
				if i == chosen {
					continue
				}
				pools = append(pools, p)
			}
			cfg.Pools = pools
			// 保存调整后的运行时配置
			_ = s.configSvc.SaveConfig(cfg)
		}
	}

	configPath := s.configSvc.GetConfigPath()
	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		return fmt.Errorf("获取配置文件路径失败: %w", err)
	}

	s.cmd = exec.Command(exePath, "--config", absConfigPath)
	s.cmd.Dir = filepath.Dir(exePath)
	// Windows下隐藏cmd窗口
	s.cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}

	// 捕获标准输出和错误输出
	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("创建stdout管道失败: %w", err)
	}

	stderr, err := s.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("创建stderr管道失败: %w", err)
	}

	if err := s.cmd.Start(); err != nil {
		return fmt.Errorf("启动XMRig失败: %w", err)
	}

	s.isRunning = true
	s.startTime = time.Now()
	s.logBuffer = make([]string, 0, s.maxLogLines)

	// 异步读取输出
	go s.readOutput(stdout, "stdout")
	go s.readOutput(stderr, "stderr")

	// 监控进程
	go s.monitorProcess()

	return nil
}

// readOutput 读取输出
func (s *XMRigService) readOutput(reader io.Reader, source string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		s.addLog(line)

		lower := strings.ToLower(line)
		if strings.Contains(lower, "new job") || strings.Contains(lower, "connected") || strings.Contains(lower, "login succeeded") {
			s.mutex.Lock()
			s.poolConnected = true
			s.mutex.Unlock()
		}
		if strings.Contains(lower, "failed") || strings.Contains(lower, "error") || strings.Contains(lower, "banned") || strings.Contains(lower, "access denied") || strings.Contains(lower, "timeout") {
			if strings.Contains(lower, "pool") || strings.Contains(lower, "net") {
				s.mutex.Lock()
				s.poolConnected = false
				s.mutex.Unlock()
			}
		}

		if s.ctx != nil {
			wailsruntime.EventsEmit(s.ctx, "miner:log", map[string]string{
				"source": source,
				"line":   line,
				"time":   time.Now().Format("15:04:05"),
			})
		}
	}
}

// addLog 添加日志
func (s *XMRigService) addLog(line string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.logBuffer) >= s.maxLogLines {
		s.logBuffer = s.logBuffer[1:]
	}
	s.logBuffer = append(s.logBuffer, line)
}

// monitorProcess 监控进程
func (s *XMRigService) monitorProcess() {
	if s.cmd != nil && s.cmd.Process != nil {
		s.cmd.Wait()
		s.mutex.Lock()
		s.isRunning = false
		s.mutex.Unlock()

		if s.ctx != nil {
			wailsruntime.EventsEmit(s.ctx, "miner:stopped", nil)
		}
	}
}

// Stop 停止挖矿
func (s *XMRigService) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 优先按PID停止当前跟踪的进程
	if s.cmd != nil && s.cmd.Process != nil {
		pid := s.cmd.Process.Pid
		s.runSilent("taskkill", "/F", "/T", "/PID", strconv.Itoa(pid))
		_ = s.cmd.Process.Kill()
	}

	s.runSilent("taskkill", "/F", "/IM", "xmrig.exe")

	s.isRunning = false
	return nil
}

// IsRunning 检查是否运行中
func (s *XMRigService) IsRunning() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.isRunning
}

// GetStatus 获取挖矿状态
func (s *XMRigService) GetStatus() (*models.MinerStatus, error) {
	s.mutex.RLock()
	running := s.isRunning
	startTime := s.startTime
	s.mutex.RUnlock()

	status := &models.MinerStatus{
		Running:   running,
		Connected: false,
	}

	if running {
		status.Uptime = int64(time.Since(startTime).Seconds())

		// 尝试从HTTP API获取详细状态
		config, err := s.configSvc.LoadConfig()
		if err == nil && config.HTTP.Enabled {
			apiStatus, err := s.getAPIStatus(config.HTTP.Host, config.HTTP.Port)
			if err == nil {
				status.Hashrate = apiStatus.Hashrate
				status.Threads = apiStatus.Threads
				if len(config.Pools) > 0 {
					status.Pool = config.Pools[0].URL
				}
			}
		}

		if err == nil && len(config.Pools) > 0 {
			s.mutex.RLock()
			connected := s.poolConnected
			s.mutex.RUnlock()
			if connected {
				status.Connected = true
			} else {
				status.Connected = s.isPoolReachable(config.Pools[0].URL)
			}
		}
	}

	return status, nil
}

func (s *XMRigService) isPoolReachable(poolURL string) bool {
	u := strings.TrimSpace(poolURL)
	if u == "" {
		return false
	}

	tlsMode := false
	if idx := strings.Index(u, "://"); idx != -1 {
		scheme := strings.ToLower(u[:idx])
		tlsMode = strings.Contains(scheme, "ssl")
		u = u[idx+3:]
	}

	addr := u
	if !strings.Contains(addr, ":") {
		addr = addr + ":443"
		if !tlsMode {
			addr = u + ":80"
		}
	}

	d := &net.Dialer{Timeout: 2 * time.Second}
	if tlsMode {
		conn, err := tls.DialWithDialer(d, "tcp", addr, &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return false
		}
		_ = conn.Close()
		return true
	}
	conn, err := d.Dial("tcp", addr)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// APIResponse XMRig API响应结构
type APIResponse struct {
	Hashrate struct {
		Total []float64 `json:"total"`
	} `json:"hashrate"`
	Resources struct {
		Threads int `json:"threads"`
	} `json:"resources"`
}

// getAPIStatus 从API获取状态
func (s *XMRigService) getAPIStatus(host string, port int) (*models.MinerStatus, error) {
	url := fmt.Sprintf("http://%s:%d/1/summary", host, port)

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	hashrate := 0.0
	if len(apiResp.Hashrate.Total) > 0 {
		hashrate = apiResp.Hashrate.Total[0]
	}

	return &models.MinerStatus{
		Hashrate: hashrate,
		Threads:  apiResp.Resources.Threads,
	}, nil
}

// GetLogs 获取日志
func (s *XMRigService) GetLogs() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	logs := make([]string, len(s.logBuffer))
	copy(logs, s.logBuffer)
	return logs
}

// ClearLogs 清空日志
func (s *XMRigService) ClearLogs() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.logBuffer = make([]string, 0, s.maxLogLines)
}

// GetSystemInfo 获取系统信息
func (s *XMRigService) GetSystemInfo() (*models.SystemInfo, error) {
	xmrigVersion := s.getXMRigVersion()

	return &models.SystemInfo{
		OS:           runtime.GOOS,
		Arch:         runtime.GOARCH,
		CPUCores:     runtime.NumCPU(),
		CPUModel:     "Unknown",
		TotalMemory:  0,
		XMRigVersion: xmrigVersion,
	}, nil
}

// getXMRigVersion 获取XMRig版本号
func (s *XMRigService) getXMRigVersion() string {
	exePath, err := s.getXMRigExecutable()
	if err != nil {
		return "Unknown"
	}

	cmd := exec.Command(exePath, "--version")
	// Windows下隐藏cmd窗口
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}
	output, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}

	// 解析版本信息，通常第一行包含版本号
	// 例如: XMRig 6.24.0
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		line := strings.TrimSpace(lines[0])
		// 提取版本号
		if strings.Contains(line, "XMRig") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1]
			}
		}
		return line
	}

	return "Unknown"
}
