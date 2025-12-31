package cmd

import (
	"github.com/sqot0/packsmith/backend/internal/config"
	"github.com/sqot0/packsmith/backend/internal/fs"
	"github.com/sqot0/packsmith/backend/internal/installer"
	"github.com/sqot0/packsmith/backend/internal/logger"
	"github.com/sqot0/packsmith/backend/internal/sources"
	"github.com/sqot0/packsmith/backend/internal/updater"
)

func (a *App) SearchMods(query string, platform string) ([]sources.ModSearch, error) {
	logger.Log.Printf("Searching mods with query: %s on platform: %s", query, platform)
	cfg, err := a.loadConfig()
	if err != nil {
		logger.Log.Printf("Failed to load config for SearchMods: %v", err)
		return nil, err
	}

	results, err := sources.SearchMods(cfg, query, platform)
	if err != nil {
		logger.Log.Printf("Error searching mods: %v", err)
		return nil, err
	}
	logger.Log.Printf("Found %d mods matching query", len(results))
	return results, nil
}

func (a *App) AddMod(modID, platform string, metadata sources.ModMetaData) error {
	logger.Log.Printf("Adding mod ID: %s from platform: %s", modID, platform)
	cfg, err := a.loadConfig()
	if err != nil {
		logger.Log.Printf("Failed to load config for AddMod: %v", err)
		return err
	}

	logger.Log.Printf("Getting download URL for mod: %s", modID)
	url, err := sources.GetDownloadURL(cfg, modID, platform, metadata.Version)
	if err != nil {
		logger.Log.Printf("Error getting download URL: %v", err)
		return err
	}
	logger.Log.Printf("Download URL obtained: %s, version: %s", url, metadata.Version)

	logger.Log.Printf("Downloading mod file")
	filename, err := fs.Download(a.ProjectPath, url, metadata.Version)
	if err != nil {
		logger.Log.Printf("Error downloading mod: %v", err)
		return err
	}
	logger.Log.Printf("Mod downloaded successfully, filename: %s", filename)

	cfg.Mods[modID] = config.Mod{
		Source:   metadata.URL,
		URL:      url,
		Version:  metadata.Version,
		Side:     metadata.Side,
		Filename: filename,
	}
	logger.Log.Printf("Mod added to config: %s", modID)

	if err := config.Save(cfg); err != nil {
		logger.Log.Printf("Error saving config: %v", err)
		return err
	}
	logger.Log.Println("Config saved successfully")
	return nil
}

func (a *App) RemoveMod(modID string) error {
	logger.Log.Printf("Removing mod ID: %s", modID)
	cfg, err := a.loadConfig()
	if err != nil {
		logger.Log.Printf("Failed to load config for RemoveMod: %v", err)
		return err
	}

	mod, ok := cfg.Mods[modID]
	if !ok {
		logger.Log.Printf("Mod %s not found in config", modID)
		return nil
	}

	logger.Log.Printf("Deleting mod file: %s", mod.Filename)
	if err := fs.Delete(a.ProjectPath, mod.Filename); err != nil {
		logger.Log.Printf("Error deleting mod file: %v", err)
		return err
	}

	delete(cfg.Mods, modID)
	logger.Log.Printf("Mod removed from config: %s", modID)

	if err := config.Save(cfg); err != nil {
		logger.Log.Printf("Error saving config: %v", err)
		return err
	}
	logger.Log.Println("Config saved successfully")
	return nil
}

func (a *App) ChangeModSide(modID, side string) error {
	logger.Log.Printf("Changing side for mod ID: %s to: %s", modID, side)
	cfg, err := a.loadConfig()
	if err != nil {
		logger.Log.Printf("Failed to load config for ChangeModSide: %v", err)
		return err
	}

	mod, ok := cfg.Mods[modID]
	if !ok {
		logger.Log.Printf("Mod %s not found in config", modID)
		return nil
	}

	mod.Side = side
	cfg.Mods[modID] = mod
	logger.Log.Printf("Mod side updated: %s", modID)

	if err := config.Save(cfg); err != nil {
		logger.Log.Printf("Error saving config: %v", err)
		return err
	}
	logger.Log.Println("Config saved successfully")
	return nil
}

func (a *App) ChangeModLocked(modID string, lock bool) error {
	logger.Log.Printf("Changing lock status for mod ID: %s to: %t", modID, lock)
	cfg, err := a.loadConfig()
	if err != nil {
		logger.Log.Printf("Failed to load config for ChangeModLocked: %v", err)
		return err
	}

	mod, ok := cfg.Mods[modID]
	if !ok {
		logger.Log.Printf("Mod %s not found in config", modID)
		return nil
	}

	mod.Locked = lock
	cfg.Mods[modID] = mod
	logger.Log.Printf("Mod lock status updated: %s", modID)

	if err := config.Save(cfg); err != nil {
		logger.Log.Printf("Error saving config: %v", err)
		return err
	}
	logger.Log.Println("Config saved successfully")
	return nil
}

func (a *App) GetModVersions(modID string) ([]string, error) {
	logger.Log.Printf("Getting versions for mod ID: %s", modID)
	cfg, err := a.loadConfig()
	if err != nil {
		logger.Log.Printf("Failed to load config for GetModVersions: %v", err)
		return nil, err
	}

	mod, ok := cfg.Mods[modID]
	if !ok {
		logger.Log.Printf("Mod %s not found in config", modID)
		return nil, nil
	}

	platform := sources.GetModPlatform(mod.Source)
	versions, err := sources.GetModVersions(cfg, modID, platform)
	if err != nil {
		logger.Log.Printf("Error getting mod versions: %v", err)
		return nil, err
	}
	logger.Log.Printf("Found %d versions for mod ID: %s", len(versions), modID)

	return versions, nil
}

func (a *App) ChangeModVersion(modID, version string) error {
	logger.Log.Printf("Changing version for mod ID: %s to: %s", modID, version)
	cfg, err := a.loadConfig()
	if err != nil {
		logger.Log.Printf("Failed to load config for ChangeModVersion: %v", err)
		return err
	}

	mod, ok := cfg.Mods[modID]
	if !ok {
		logger.Log.Printf("Mod %s not found in config", modID)
		return nil
	}

	url, err := sources.GetDownloadURL(cfg, modID, sources.GetModPlatform(mod.Source), version)
	if err != nil {
		logger.Log.Printf("Error getting download URL: %v", err)
		return err
	}
	logger.Log.Printf("Download URL obtained: %s, version: %s", url, version)

	logger.Log.Printf("Downloading mod file")
	filename, err := fs.Download(a.ProjectPath, url, version)

	mod.Version = version
	mod.URL = url
	mod.Filename = filename
	cfg.Mods[modID] = mod
	logger.Log.Printf("Mod version updated: %s", modID)

	if err := config.Save(cfg); err != nil {
		logger.Log.Printf("Error saving config: %v", err)
		return err
	}

	return nil
}

func (a *App) CheckModsUpdates(modIDs []string) ([]updater.ModToUpdate, error) {
	logger.Log.Printf("Checking updates for %d mods", len(modIDs))
	cfg, err := a.loadConfig()
	if err != nil {
		logger.Log.Printf("Failed to load config for CheckModsUpdates: %v", err)
		return nil, err
	}

	modsToUpdate, err := updater.CheckMods(cfg, modIDs)
	if err != nil {
		logger.Log.Printf("Error checking mods updates: %v", err)
		return nil, err
	}
	logger.Log.Printf("Found %d mods to update", len(modsToUpdate))
	return modsToUpdate, nil
}

func (a *App) UpdateMods(modsToUpdate []updater.ModToUpdate) error {
	logger.Log.Printf("Updating %d mods", len(modsToUpdate))
	cfg, err := a.loadConfig()
	if err != nil {
		logger.Log.Printf("Failed to load config for UpdateMods: %v", err)
		return err
	}

	if err := updater.UpdateMods(cfg, modsToUpdate, a.ProjectPath); err != nil {
		logger.Log.Printf("Error updating mods: %v", err)
		return err
	}
	logger.Log.Println("Mods updated successfully")
	return nil
}

func (a *App) InstallMods() error {
	logger.Log.Println("Installing mods")
	err := installer.InstallMods(a.ProjectPath)
	if err != nil {
		logger.Log.Printf("Error installing mods: %v", err)
		return err
	}
	logger.Log.Println("Mods installed successfully")
	return nil
}
