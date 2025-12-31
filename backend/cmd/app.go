package cmd

import (
	"context"
	"os"

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
}

func (a *App) GetLogs() (string, error) {
	data, err := os.ReadFile("app.log")
	if err != nil {
		logger.Log.Printf("Error reading log file: %v", err)
		return "", err
	}
	return string(data), nil
}

func (a *App) loadConfig() (*config.Config, error) {
	cfg, err := config.Load(a.ProjectPath)
	if err != nil {
		logger.Log.Printf("Error loading config: %v", err)
		return nil, err
	}
	return cfg, nil
}
