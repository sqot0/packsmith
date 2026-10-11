package cmd

import (
	"context"
	"os"

	"github.com/sqot0/packsmith/backend/internal/appconfig"
	"github.com/sqot0/packsmith/backend/internal/config"
	"github.com/sqot0/packsmith/backend/internal/discord"
	"github.com/sqot0/packsmith/backend/internal/logger"
)

type App struct {
	ctx         context.Context
	ProjectPath string
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	logger.Init()
	discord.Init()
	logger.Log.Println("Logger initialized successfully")

	if _, err := appconfig.Load(); err != nil {
		logger.Log.Printf("Failed to load app config: %v", err)
	}
}

func (a *App) GetLogs() (string, error) {
	data, err := os.ReadFile("app.log")
	if err != nil {
		logger.Log.Printf("Error reading log file: %v", err)
		return "", err
	}
	return string(data), nil
}

func (a *App) GetProjectConfig() (*config.Config, error) {
	return a.loadConfig()
}

func (a *App) SetProjectPath(path string) {
	a.ProjectPath = path
}

func (a *App) GetAppSettings() (*appconfig.AppSettings, error) {
	return appconfig.Load()
}

func (a *App) SaveAppSettings(settings *appconfig.AppSettings) error {
	logger.Log.Printf("Saving app settings: maxWorkers=%d", settings.MaxWorkers)
	return appconfig.Save(settings)
}

func (a *App) GetRecentProjects() ([]appconfig.RecentProject, error) {
	cfg, err := appconfig.Load()
	if err != nil {
		return nil, err
	}
	return cfg.RecentProjects, nil
}

func (a *App) ClearRecentProjects() error {
	logger.Log.Println("Clearing recent projects list")
	_, err := appconfig.ClearRecentProjects()
	return err
}

func (a *App) RemoveRecentProject(path string) error {
	logger.Log.Printf("Removing recent project: %s", path)
	_, err := appconfig.RemoveRecentProject(path)
	return err
}

func (a *App) loadConfig() (*config.Config, error) {
	cfg, err := config.Load(a.ProjectPath)
	if err != nil {
		logger.Log.Printf("Error loading config: %v", err)
		return nil, err
	}
	return cfg, nil
}
