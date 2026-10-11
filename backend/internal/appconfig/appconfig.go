package appconfig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type RecentProject struct {
	Path       string    `json:"path"`
	Name       string    `json:"name"`
	LastOpened time.Time `json:"lastOpened"`
}

type AppSettings struct {
	MaxWorkers     int             `json:"maxWorkers"`
	RecentProjects []RecentProject `json:"recentProjects"`
}

const (
	DefaultMaxWorkers = 15
	ConfigDirName     = "packsmith"
	ConfigFileName    = "config.json"
	MaxRecentProjects = 10
)

var (
	mu          sync.RWMutex
	cachedLimit int = DefaultMaxWorkers
)

func GetConfigFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config dir: %w", err)
	}

	appDir := filepath.Join(configDir, ConfigDirName)
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create config dir: %w", err)
	}

	return filepath.Join(appDir, ConfigFileName), nil
}

func Load() (*AppSettings, error) {
	mu.Lock()
	defer mu.Unlock()

	configPath, err := GetConfigFilePath()
	if err != nil {
		return &AppSettings{MaxWorkers: cachedLimit, RecentProjects: []RecentProject{}}, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			defaultCfg := &AppSettings{
				MaxWorkers:     DefaultMaxWorkers,
				RecentProjects: make([]RecentProject, 0),
			}
			cachedLimit = defaultCfg.MaxWorkers
			_ = saveLocked(configPath, defaultCfg)
			return defaultCfg, nil
		}
		return &AppSettings{MaxWorkers: cachedLimit, RecentProjects: []RecentProject{}}, err
	}

	trimmed := bytes.TrimRight(data, " \t\r\n\x00")

	var cfg AppSettings
	if err := json.Unmarshal(trimmed, &cfg); err != nil {
		defaultCfg := &AppSettings{
			MaxWorkers:     DefaultMaxWorkers,
			RecentProjects: make([]RecentProject, 0),
		}
		cachedLimit = defaultCfg.MaxWorkers
		return defaultCfg, nil
	}

	if cfg.MaxWorkers <= 0 {
		cfg.MaxWorkers = DefaultMaxWorkers
	}
	if cfg.RecentProjects == nil {
		cfg.RecentProjects = make([]RecentProject, 0)
	}

	cachedLimit = cfg.MaxWorkers
	return &cfg, nil
}

func Save(cfg *AppSettings) error {
	mu.Lock()
	defer mu.Unlock()

	configPath, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	if cfg.MaxWorkers <= 0 {
		cfg.MaxWorkers = DefaultMaxWorkers
	}
	if cfg.RecentProjects == nil {
		cfg.RecentProjects = make([]RecentProject, 0)
	}

	cachedLimit = cfg.MaxWorkers
	return saveLocked(configPath, cfg)
}

func saveLocked(configPath string, cfg *AppSettings) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}
	data = append(data, '\n')

	f, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	_ = f.Truncate(0)
	n, err := f.Write(data)
	if err != nil {
		return err
	}
	_ = f.Truncate(int64(n))
	_ = f.Sync()
	return nil
}

func GetMaxWorkers() int {
	mu.RLock()
	workers := cachedLimit
	mu.RUnlock()

	if workers > 0 {
		return workers
	}

	cfg, err := Load()
	if err == nil && cfg.MaxWorkers > 0 {
		return cfg.MaxWorkers
	}

	return DefaultMaxWorkers
}

func AddRecentProject(path string, name string) (*AppSettings, error) {
	if strings.TrimSpace(path) == "" {
		return Load()
	}

	cleanPath := filepath.Clean(path)
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	mu.Lock()
	defer mu.Unlock()

	configPath, err := GetConfigFilePath()
	if err != nil {
		return nil, err
	}

	updated := make([]RecentProject, 0, len(cfg.RecentProjects)+1)
	updated = append(updated, RecentProject{
		Path:       cleanPath,
		Name:       name,
		LastOpened: time.Now(),
	})

	for _, p := range cfg.RecentProjects {
		if strings.EqualFold(filepath.Clean(p.Path), cleanPath) {
			continue
		}
		updated = append(updated, p)
		if len(updated) >= MaxRecentProjects {
			break
		}
	}

	cfg.RecentProjects = updated
	if err := saveLocked(configPath, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func RemoveRecentProject(path string) (*AppSettings, error) {
	cleanPath := filepath.Clean(path)
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	mu.Lock()
	defer mu.Unlock()

	configPath, err := GetConfigFilePath()
	if err != nil {
		return nil, err
	}

	updated := make([]RecentProject, 0, len(cfg.RecentProjects))
	for _, p := range cfg.RecentProjects {
		if strings.EqualFold(filepath.Clean(p.Path), cleanPath) {
			continue
		}
		updated = append(updated, p)
	}

	cfg.RecentProjects = updated
	if err := saveLocked(configPath, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func ClearRecentProjects() (*AppSettings, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	mu.Lock()
	defer mu.Unlock()

	configPath, err := GetConfigFilePath()
	if err != nil {
		return nil, err
	}

	cfg.RecentProjects = make([]RecentProject, 0)
	if err := saveLocked(configPath, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
