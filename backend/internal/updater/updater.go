package updater

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/sqot0/packsmith/backend/internal/config"
	"github.com/sqot0/packsmith/backend/internal/fs"
	"github.com/sqot0/packsmith/backend/internal/logger"
	"github.com/sqot0/packsmith/backend/internal/sources"
	"github.com/sqot0/packsmith/backend/internal/util"
)

type ModToUpdate struct {
	ModId   string
	Version string
	URL     string
}

func CheckMods(cfg *config.Config, modIDs []string) ([]ModToUpdate, error) {
	logger.Log.Printf("Checking updates for %d mods", len(modIDs))
	type job struct {
		modId string
		mod   config.Mod
	}

	type result struct {
		modId   string
		version string
		url     string
	}

	processJob := func(j job) result {
		logger.Log.Printf("Checking update for mod: %s", j.modId)
		platform := sources.GetModPlatform(j.mod.Source)

		version, err := sources.GetLatestVersion(cfg, j.modId, platform)
		if err != nil {
			logger.Log.Printf("Error checking mod %s: %v", j.modId, err)
			return result{modId: j.modId, version: version, url: j.mod.URL}
		}
		logger.Log.Printf("Mod %s latest version: %s", j.modId, version)

		url, err := sources.GetDownloadURL(cfg, j.modId, platform, version)
		return result{
			modId:   j.modId,
			version: version,
			url:     url,
		}
	}

	jobs := make(chan job)
	results := util.WorkerPool(jobs, processJob, len(modIDs))

	go func() {
		for _, modID := range modIDs {
			mod, ok := cfg.Mods[modID]
			if !ok || mod.Source == "" || mod.Locked {
				logger.Log.Printf("Skipping mod %s (not found, no source, or locked)", modID)
				continue
			}
			jobs <- job{modId: modID, mod: mod}
		}
		close(jobs)
	}()

	modsToUpdate := make([]ModToUpdate, 0)

	for r := range results {
		if r.version != "" && r.version != cfg.Mods[r.modId].Version {
			logger.Log.Printf("Mod %s needs update from %s to %s", r.modId, cfg.Mods[r.modId].Version, r.version)
			modsToUpdate = append(modsToUpdate, ModToUpdate{
				ModId:   r.modId,
				Version: r.version,
				URL:     r.url,
			})
		} else {
			logger.Log.Printf("Mod %s is up to date", r.modId)
		}
	}

	logger.Log.Printf("Found %d mods to update", len(modsToUpdate))
	return modsToUpdate, nil
}

func UpdateMods(cfg *config.Config, mods []ModToUpdate, projectPath string) error {
	logger.Log.Printf("Updating %d mods", len(mods))
	var mx sync.Mutex

	processUpdate := func(mod ModToUpdate) error {
		logger.Log.Printf("Updating mod: %s to version: %s", mod.ModId, mod.Version)
		filename, err := fs.Download(projectPath, mod.URL, mod.Version)
		if err != nil {
			logger.Log.Printf("Error downloading mod %s: %v", mod.ModId, err)
			return fmt.Errorf("%s: %w", mod.ModId, err)
		}

		modCfg := cfg.Mods[mod.ModId]

		logger.Log.Printf("Removing old cache file for mod: %s", mod.ModId)
		if err := os.RemoveAll(filepath.Join(projectPath, "cache", modCfg.Filename)); err != nil {
			logger.Log.Printf("Error removing old cache file for %s: %v", mod.ModId, err)
			return fmt.Errorf("%s: %w", mod.ModId, err)
		}

		mx.Lock()
		modCfg.Version = mod.Version
		modCfg.Filename = filename
		modCfg.URL = mod.URL
		cfg.Mods[mod.ModId] = modCfg
		mx.Unlock()
		logger.Log.Printf("Mod %s updated successfully", mod.ModId)
		return nil
	}

	jobs := make(chan ModToUpdate, len(mods))
	results := util.WorkerPool(jobs, processUpdate, len(mods))

	go func() {
		for _, mod := range mods {
			jobs <- mod
		}
		close(jobs)
	}()

	for err := range results {
		if err != nil {
			logger.Log.Printf("Error updating mod: %v", err)
		}
	}

	logger.Log.Println("Saving updated config")
	err := config.Save(cfg)
	if err != nil {
		logger.Log.Printf("Error saving config: %v", err)
		return err
	}
	logger.Log.Println("Config saved successfully")
	return nil
}
