package installer

import (
	"os"
	"path"

	"github.com/sqot0/packsmith/backend/internal/config"
	"github.com/sqot0/packsmith/backend/internal/fs"
	"github.com/sqot0/packsmith/backend/internal/logger"
	"github.com/sqot0/packsmith/backend/internal/util"
)

func InstallMods(projectPath string) error {
	logger.Log.Printf("Installing mods for project: %s", projectPath)
	cfg, err := config.Load(projectPath)
	if err != nil {
		logger.Log.Printf("Error loading config: %v", err)
		return err
	}

	cacheFolder := path.Join(projectPath, "cache")
	clientFolder := path.Join(projectPath, "client")
	serverFolder := path.Join(projectPath, "server")

	logger.Log.Println("Cleaning and creating client and server folders")
	for _, folder := range []string{clientFolder, serverFolder} {
		if err := os.RemoveAll(folder); err != nil {
			logger.Log.Printf("Error removing folder %s: %v", folder, err)
			return err
		}
		if err := os.MkdirAll(folder, 0o755); err != nil {
			logger.Log.Printf("Error creating folder %s: %v", folder, err)
			return err
		}
	}

	processMod := func(mod config.Mod) error {
		logger.Log.Printf("Processing mod: %s", mod.Filename)
		cacheMod := path.Join(cacheFolder, mod.Filename)

		if _, err := os.Stat(cacheMod); os.IsNotExist(err) {
			logger.Log.Printf("Mod not in cache, downloading: %s", mod.URL)
			if _, err := fs.Download(projectPath, mod.URL, mod.Version); err != nil {
				logger.Log.Printf("Error downloading mod: %v", err)
				return err
			}
		}

		switch mod.Side {
		case "both":
			logger.Log.Printf("Copying mod to client and server: %s", mod.Filename)
			if err := fs.Copy(cacheMod, path.Join(clientFolder, mod.Filename)); err != nil {
				logger.Log.Printf("Error copying to client: %v", err)
				return err
			}
			if err := fs.Copy(cacheMod, path.Join(serverFolder, mod.Filename)); err != nil {
				logger.Log.Printf("Error copying to server: %v", err)
				return err
			}

		case "client":
			logger.Log.Printf("Copying mod to client: %s", mod.Filename)
			if err := fs.Copy(cacheMod, path.Join(clientFolder, mod.Filename)); err != nil {
				logger.Log.Printf("Error copying to client: %v", err)
				return err
			}

		case "server":
			logger.Log.Printf("Copying mod to server: %s", mod.Filename)
			if err := fs.Copy(cacheMod, path.Join(serverFolder, mod.Filename)); err != nil {
				logger.Log.Printf("Error copying to server: %v", err)
				return err
			}
		}

		return nil
	}

	jobs := make(chan config.Mod, len(cfg.Mods))
	results := util.WorkerPool(jobs, processMod, len(cfg.Mods))

	go func() {
		for _, mod := range cfg.Mods {
			jobs <- mod
		}
		close(jobs)
	}()

	for err := range results {
		if err != nil {
			logger.Log.Printf("Error processing mod: %v", err)
			return err
		}
	}

	logger.Log.Println("Mods installed successfully")
	return nil
}
