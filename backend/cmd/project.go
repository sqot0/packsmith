package cmd

import (
	"github.com/sqot0/packsmith/backend/internal/config"
	"github.com/sqot0/packsmith/backend/internal/discord"
	"github.com/sqot0/packsmith/backend/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) SelectProjectDirectory() (string, error) {
	logger.Log.Println("Opening directory selection dialog")
	projectPath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:                "Open Modpack Project",
		CanCreateDirectories: true,
	})
	if err != nil {
		logger.Log.Printf("Error opening directory dialog: %v", err)
		return "", err
	}
	logger.Log.Printf("Selected project directory: %s", projectPath)
	return projectPath, nil
}

func (a *App) OpenProject(projectPath string) (*config.Config, error) {
	logger.Log.Printf("Opening project at path: %s", projectPath)
	cfg, err := config.Load(projectPath)
	if err != nil {
		logger.Log.Printf("Error loading project config: %v", err)
		return nil, err
	}
	a.ProjectPath = projectPath
	discord.OpenProject(cfg)
	logger.Log.Println("Project opened successfully")
	return cfg, nil
}

func (a *App) InitializeProject(projectPath, name, mc, loader string) error {
	logger.Log.Printf("Initializing project: %s with MC version: %s and with loader: %s", name, mc, loader)
	if err := config.Init(projectPath, name, mc, loader); err != nil {
		logger.Log.Printf("Error initializing project: %v", err)
		return err
	}
	logger.Log.Println("Project initialized successfully")
	return nil
}
