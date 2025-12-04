package api

import (
	"fmt"
	"go-wails/internal/models"
	"go-wails/internal/service"
)

// MinerAPI 挖矿API
type MinerAPI struct {
	xmrigService  *service.XMRigService
	configService *service.ConfigService
}

// NewMinerAPI 创建挖矿API
func NewMinerAPI(xmrigService *service.XMRigService, configService *service.ConfigService) *MinerAPI {
	return &MinerAPI{
		xmrigService:  xmrigService,
		configService: configService,
	}
}

// StartMining 开始挖矿
func (api *MinerAPI) StartMining() error {
	return api.xmrigService.Start()
}

// StopMining 停止挖矿
func (api *MinerAPI) StopMining() error {
	return api.xmrigService.Stop()
}

// GetMinerStatus 获取挖矿状态
func (api *MinerAPI) GetMinerStatus() (*models.MinerStatus, error) {
	return api.xmrigService.GetStatus()
}

// GetSystemInfo 获取系统信息
func (api *MinerAPI) GetSystemInfo() (*models.SystemInfo, error) {
	return api.xmrigService.GetSystemInfo()
}

// GetLogs 获取日志
func (api *MinerAPI) GetLogs() []string {
	return api.xmrigService.GetLogs()
}

// ClearLogs 清空日志
func (api *MinerAPI) ClearLogs() {
	api.xmrigService.ClearLogs()
}

// LoadConfig 加载配置
func (api *MinerAPI) LoadConfig() (*models.XMRigConfig, error) {
	return api.configService.LoadConfig()
}

// SaveConfig 保存配置
func (api *MinerAPI) SaveConfig(config *models.XMRigConfig) error {
	if api.xmrigService.IsRunning() {
		return fmt.Errorf("挖矿运行中，禁止修改配置")
	}
	return api.configService.SaveConfig(config)
}

// GetDefaultConfig 获取默认配置
func (api *MinerAPI) GetDefaultConfig() *models.XMRigConfig {
	return api.configService.GetDefaultConfig()
}
