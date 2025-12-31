package sources

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sqot0/packsmith/backend/internal/config"
	"github.com/sqot0/packsmith/backend/internal/logger"
)

func setHeadersForRequest(req *http.Request) {
	logger.Log.Println("Setting headers for HTTP request")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
}

type ModSearch struct {
	ID          string
	Name        string
	Description string
	ClientSide  string
	ServerSide  string
	Downloads   string
	URL         string
	Versions    []string
}

type ModMetaData struct {
	URL, Side, Version string
}

func SearchMods(cfg *config.Config, query, platform string) ([]ModSearch, error) {
	logger.Log.Printf("Searching mods on platform: %s with query: %s", platform, query)
	switch platform {
	case "modrinth":
		return searchModsModrinth(cfg, query)
	case "curseforge":
		return searchModsCurseforge(cfg, query)
	default:
		logger.Log.Printf("Unknown platform: %s", platform)
		return nil, fmt.Errorf("unknown platform: %s", platform)
	}
}

func GetDownloadURL(cfg *config.Config, modID, platform, version string) (string, error) {
	logger.Log.Printf("Getting download URL for mod: %s on platform: %s", modID, platform)
	switch platform {
	case "modrinth":
		return getDownloadURLModrinth(cfg, modID, version)
	case "curseforge":
		return getDownloadURLCurseforge(cfg, modID, version)
	default:
		logger.Log.Printf("Unknown platform: %s", platform)
		return "", fmt.Errorf("unknown platform: %s", platform)
	}
}

func GetModVersions(cfg *config.Config, modID, platform string) ([]string, error) {
	logger.Log.Printf("Getting versions for mod: %s on platform: %s", modID, platform)
	switch platform {
	case "modrinth":
		return getModVersionsModrinth(cfg, modID)
	case "curseforge":
		return getModVersionsCurseforge(cfg, modID)
	default:
		logger.Log.Printf("Unknown platform: %s", platform)
		return nil, fmt.Errorf("unknown platform: %s", platform)
	}
}

func GetLatestVersion(cfg *config.Config, modID, platform string) (string, error) {
	logger.Log.Printf("Getting latest version for mod: %s on platform: %s", modID, platform)
	switch platform {
	case "modrinth":
		return getLatestVersionModrinth(cfg, modID)
	case "curseforge":
		return getLatestVersionCurseforge(cfg, modID)
	default:
		logger.Log.Printf("Unknown platform: %s", platform)
		return "", fmt.Errorf("unknown platform: %s", platform)
	}
}

func GetModPlatform(source string) string {
	logger.Log.Printf("Determining platform for source: %s", source)
	if strings.HasPrefix(source, "https://www.curseforge.com") {
		logger.Log.Println("Platform: curseforge")
		return "curseforge"
	}
	logger.Log.Println("Platform: modrinth")
	return "modrinth"
}
