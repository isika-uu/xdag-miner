package service

import (
	"embed"
	"encoding/json"
	"fmt"
	"go-wails/internal/models"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed xmrig-embedded/*
var embeddedConfig embed.FS

// ConfigService 配置服务
type ConfigService struct {
	configPath string
	runtimeDir string
}

// NewConfigService 创建配置服务
func NewConfigService() *ConfigService {
	runtimeDir := filepath.Join(os.TempDir(), "xmrig-runtime")
	return &ConfigService{
		runtimeDir: runtimeDir,
	}
}

// GetConfigPath 获取配置文件路径
func (s *ConfigService) GetConfigPath() string {
	os.MkdirAll(s.runtimeDir, 0755)
	return filepath.Join(s.runtimeDir, "config.json")
}

// ensureConfigExists 确保配置文件存在
func (s *ConfigService) ensureConfigExists() error {
	configPath := s.GetConfigPath()

	if _, err := os.Stat(configPath); err == nil {
		return nil
	}

	arch := runtime.GOARCH
	var embeddedPath string
	if arch == "arm64" {
		embeddedPath = "xmrig-embedded/xmrig-windows-arm64/config.json"
	} else {
		embeddedPath = "xmrig-embedded/xmrig-windows-amd64/config.json"
	}

	data, err := embeddedConfig.ReadFile(embeddedPath)
	if err != nil {
		defaultConfig := s.GetDefaultConfig()
		return s.SaveConfig(defaultConfig)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("创建配置文件失败: %w", err)
	}

	return nil
}

// LoadConfig 加载配置
func (s *ConfigService) LoadConfig() (*models.XMRigConfig, error) {
	if err := s.ensureConfigExists(); err != nil {
		return nil, err
	}

	configPath := s.GetConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config models.XMRigConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &config, nil
}

// SaveConfig 保存配置
func (s *ConfigService) SaveConfig(config *models.XMRigConfig) error {
	configPath := s.GetConfigPath()

	// 读取现有配置为通用Map，保留未知字段；若不存在则直接写入新配置
	existingData, err := os.ReadFile(configPath)
	var existing map[string]interface{}
	if err == nil {
		if err := json.Unmarshal(existingData, &existing); err != nil {
			return fmt.Errorf("解析现有配置失败: %w", err)
		}
	}

	// 将更新配置转为Map
	updateBytes, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化更新配置失败: %w", err)
	}
	var updates map[string]interface{}
	if err := json.Unmarshal(updateBytes, &updates); err != nil {
		return fmt.Errorf("解析更新配置失败: %w", err)
	}

	// 递归合并：仅覆盖传入的键；保留未知键
	var finalBytes []byte
	if existing != nil {
		merged := mergeJSON(existing, updates)
		finalBytes, err = json.MarshalIndent(merged, "", "    ")
	} else {
		// 无现有配置，直接写入更新配置
		finalBytes, err = json.MarshalIndent(updates, "", "    ")
	}
	if err != nil {
		return fmt.Errorf("序列化合并后的配置失败: %w", err)
	}
	if err := os.WriteFile(configPath, finalBytes, 0644); err != nil {
		return fmt.Errorf("保存配置文件失败: %w", err)
	}
	return nil
}

// mergeJSON 递归合并对象：map合并，数组整体替换，原子值覆盖
func mergeJSON(dst, src map[string]interface{}) map[string]interface{} {
	if dst == nil {
		dst = map[string]interface{}{}
	}
	for k, v := range src {
		switch nv := v.(type) {
		case map[string]interface{}:
			if dv, ok := dst[k].(map[string]interface{}); ok {
				dst[k] = mergeJSON(dv, nv)
			} else {
				// 用新的对象替换
				dst[k] = mergeJSON(map[string]interface{}{}, nv)
			}
		case []interface{}:
			// 切片整体替换（如 pools 等）
			dst[k] = nv
		default:
			// 原子类型覆盖，包含 null
			dst[k] = nv
		}
	}
	return dst
}

// GetDefaultConfig 获取默认配置
func (s *ConfigService) GetDefaultConfig() *models.XMRigConfig {
	pass := "x"
	return &models.XMRigConfig{
		Autosave: true,
		CPU: models.CPUConfig{
			Enabled:        true,
			HugePages:      true,
			MaxThreadsHint: 100,
			ASM:            true,
		},
		Pools: []models.PoolConfig{
			{
				URL:     "stratum+ssl://equal.xdagminer.com:13003",
				User:    "NNZabJQEhrQGTPabqABWVK9v3rSsNQ7Sy",
				Pass:    pass,
				Enabled: true,
			},
		},
		RandomX: models.RandomXConfig{
			Init:     -1,
			InitAVX2: -1,
			Mode:     "auto",
			NUMA:     true,
		},
		HTTP: models.HTTPConfig{
			Enabled:    true,
			Host:       "127.0.0.1",
			Port:       3649,
			Restricted: true,
		},
	}
}
