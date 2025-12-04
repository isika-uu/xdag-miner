package main

import (
	"context"
	"go-wails/internal/api"
	"go-wails/internal/models"
	"go-wails/internal/service"
)

// App struct
type App struct {
	ctx           context.Context
	xmrigService  *service.XMRigService
	configService *service.ConfigService
	minerAPI      *api.MinerAPI
}

// NewApp creates a new App application struct
func NewApp() *App {
	xmrigService := service.NewXMRigService()
	configService := service.NewConfigService()

	return &App{
		xmrigService:  xmrigService,
		configService: configService,
		minerAPI:      api.NewMinerAPI(xmrigService, configService),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.xmrigService.SetContext(ctx)
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	a.xmrigService.Stop()
}

// === 挖矿控制相关方法 ===

// StartMining 开始挖矿
func (a *App) StartMining() error {
	return a.minerAPI.StartMining()
}

// StopMining 停止挖矿
func (a *App) StopMining() error {
	return a.minerAPI.StopMining()
}

// GetMinerStatus 获取挖矿状态
func (a *App) GetMinerStatus() (*models.MinerStatus, error) {
	return a.minerAPI.GetMinerStatus()
}

// GetSystemInfo 获取系统信息
func (a *App) GetSystemInfo() (*models.SystemInfo, error) {
	return a.minerAPI.GetSystemInfo()
}

// GetLogs 获取日志
func (a *App) GetLogs() []string {
	return a.minerAPI.GetLogs()
}

// ClearLogs 清空日志
func (a *App) ClearLogs() {
	a.minerAPI.ClearLogs()
}

// === 配置管理相关方法 ===

// LoadConfig 加载配置
func (a *App) LoadConfig() (*models.XMRigConfig, error) {
	return a.minerAPI.LoadConfig()
}

// SaveConfig 保存配置
func (a *App) SaveConfig(config *models.XMRigConfig) error {
	return a.minerAPI.SaveConfig(config)
}

// GetDefaultConfig 获取默认配置
func (a *App) GetDefaultConfig() *models.XMRigConfig {
	return a.minerAPI.GetDefaultConfig()
}
