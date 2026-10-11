package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

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

	if loader != "forge" && loader != "fabric" && loader != "neoforge" && loader != "quilt" {
		logger.Log.Printf("Invalid loader specified: %s", loader)
		return errors.New("loader must be either 'forge', 'neoforge', 'quilt' or 'fabric'")
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
	cfgFile := filepath.Join(projectPath, "packsmith.json")
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		logger.Log.Printf("Error reading config file: %v", err)
		return nil, errors.New("initialize project before using other commands")
	}

	trimmed := bytes.TrimRight(data, " \t\r\n\x00")

	var cfg Config
	cfg.path = projectPath
	err = json.Unmarshal(trimmed, &cfg)
	if err != nil {
		logger.Log.Printf("Error unmarshaling config: %v", err)
		return nil, err
	}

	if len(trimmed) < len(data) && bytes.ContainsAny(data[len(trimmed):], " \t\x00") {
		logger.Log.Printf("Detected and cleaned trailing whitespace/null bytes from %s", cfgFile)
		_ = write(cfg)
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
	cfgFile := filepath.Join(cfg.path, "packsmith.json")
	logger.Log.Printf("Writing config to file: %s", cfgFile)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		logger.Log.Printf("Error marshaling config: %v", err)
		return err
	}
	data = append(data, '\n')

	tmpFile := filepath.Join(cfg.path, "packsmith.json.tmp")
	writtenWithTemp := false

	if f, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644); err == nil {
		n, writeErr := f.Write(data)
		_ = f.Truncate(int64(n))
		_ = f.Sync()
		closeErr := f.Close()

		if writeErr == nil && closeErr == nil && n == len(data) {
			if err := os.Rename(tmpFile, cfgFile); err == nil {
				writtenWithTemp = true
			} else {
				_ = os.Remove(cfgFile)
				if err := os.Rename(tmpFile, cfgFile); err == nil {
					writtenWithTemp = true
				}
			}
		}

		if !writtenWithTemp {
			_ = os.Remove(tmpFile)
		}
	}

	if !writtenWithTemp {
		f, err := os.OpenFile(cfgFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
		if err != nil {
			logger.Log.Printf("Error opening config file: %v", err)
			return err
		}
		defer f.Close()

		_ = f.Truncate(0)
		n, err := f.Write(data)
		if err != nil {
			logger.Log.Printf("Error writing config file: %v", err)
			return err
		}
		if err := f.Truncate(int64(n)); err != nil {
			logger.Log.Printf("Warning: failed to truncate config file to %d bytes: %v", n, err)
		}
		if err := f.Sync(); err != nil {
			logger.Log.Printf("Warning: failed to sync config file: %v", err)
		}
	}

	logger.Log.Println("Config file written successfully")
	return nil
}
