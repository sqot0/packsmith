package config

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/sqot0/packsmith/backend/internal/logger"
)

type Mod struct {
	Source   string `json:"source"`
	Side     string `json:"side"`
	Version  string `json:"version"`
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Locked   bool   `json:"locked"`
}

type Config struct {
	Name      string         `json:"name"`
	Minecraft string         `json:"minecraft"`
	Loader    string         `json:"loader"`
	Mods      map[string]Mod `json:"mods"`
	path      string
}

func Init(projectPath, name, mc, loader string) error {
	logger.Log.Printf("Initializing config for project: %s, MC: %s, Loader: %s", name, mc, loader)

	if loader != "forge" && loader != "fabric" {
		logger.Log.Printf("Invalid loader specified: %s", loader)
		return errors.New("loader must be either 'forge' or 'fabric'")
	}

	cfg := Config{name, mc, loader, map[string]Mod{}, projectPath}
	err := write(cfg)
	if err != nil {
		logger.Log.Printf("Error writing initial config: %v", err)
		return err
	}
	logger.Log.Println("Config initialized successfully")
	return nil
}

func Load(projectPath string) (*Config, error) {
	logger.Log.Printf("Loading config from path: %s", projectPath)
	cfgFile := path.Join(projectPath, "packsmith.json")
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		logger.Log.Printf("Error reading config file: %v", err)
		return nil, errors.New("initialize project before using other commands")
	}

	var cfg Config
	cfg.path = projectPath
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		logger.Log.Printf("Error unmarshaling config: %v", err)
		return nil, err
	}
	logger.Log.Println("Config loaded successfully")
	return &cfg, nil
}

func Save(cfg *Config) error {
	logger.Log.Println("Saving config")
	err := write(*cfg)
	if err != nil {
		logger.Log.Printf("Error saving config: %v", err)
		return err
	}
	logger.Log.Println("Config saved successfully")
	return nil
}

func write(cfg Config) error {
	logger.Log.Printf("Writing config to file: %s", path.Join(cfg.path, "packsmith.json"))
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		logger.Log.Printf("Error marshaling config: %v", err)
		return err
	}
	cfgFile := path.Join(cfg.path, "packsmith.json")
	err = os.WriteFile(cfgFile, data, 0o644)
	if err != nil {
		logger.Log.Printf("Error writing config file: %v", err)
		return err
	}
	logger.Log.Println("Config file written successfully")
	return nil
}
