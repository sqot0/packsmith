package sources

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sqot0/packsmith/backend/internal/config"
	"github.com/sqot0/packsmith/backend/internal/logger"
	"github.com/sqot0/packsmith/backend/internal/util"
)

const (
	curseforgeAPI = "https://api.curseforge.com/v1"
	minecraftGame = 432
)

type curseforgeSearchResponse struct {
	Data []curseforgeMod `json:"data"`
}

type curseforgeFilesResponse struct {
	Data []curseforgeFile `json:"data"`
}

type curseforgeMod struct {
	ID            int              `json:"id"`
	Name          string           `json:"name"`
	Summary       string           `json:"summary"`
	DownloadCount float64          `json:"downloadCount"`
	Slug          string           `json:"slug"`
	LatestFiles   []curseforgeFile `json:"latestFiles"`
	Links         struct {
		WebsiteURL string `json:"websiteUrl"`
	} `json:"links"`
}

type curseforgeFile struct {
	ID           int      `json:"id"`
	DisplayName  string   `json:"displayName"`
	FileName     string   `json:"fileName"`
	DownloadURL  string   `json:"downloadUrl"`
	GameVersions []string `json:"gameVersions"`
}

func curseforgeLoaderType(loader string) int {
	switch strings.ToLower(loader) {
	case "forge":
		return 1
	case "fabric":
		return 4
	case "quilt":
		return 5
	case "neoforge":
		return 6
	default:
		return 0
	}
}

func makeCurseforgeRequest(endpoint string, target interface{}) error {
	req, err := http.NewRequest("GET", curseforgeAPI+endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", "$2a$10$2nOZBkW3JY/w.42wu7fH4OiyObiV4osugysUcaw/YiGElB8dTLNfi")
	req.Header.Set("User-Agent", "PackSmith/1.0")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("curseforge api returned status %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func searchModsCurseforge(cfg *config.Config, query string) ([]ModSearch, error) {
	logger.Log.Printf("Searching CurseForge for query: %s", query)

	params := url.Values{}
	params.Set("gameId", strconv.Itoa(minecraftGame))
	params.Set("searchFilter", query)
	params.Set("pageSize", "5")
	params.Set("sortField", "2")
	params.Set("sortOrder", "desc")
	params.Set("classId", "6")

	if cfg.Minecraft != "" {
		params.Set("gameVersion", cfg.Minecraft)
	}

	modLoaderType := curseforgeLoaderType(cfg.Loader)
	if modLoaderType != 0 {
		params.Set("modLoaderType", strconv.Itoa(modLoaderType))
	}

	var response curseforgeSearchResponse

	err := makeCurseforgeRequest(
		"/mods/search?"+params.Encode(),
		&response,
	)
	if err != nil {
		logger.Log.Printf("CurseForge search error: %v", err)
		return nil, err
	}

	var (
		mods []ModSearch
		mu   sync.Mutex
	)

	processMod := func(mod curseforgeMod) error {
		versions, err := getModVersionsCurseforge(cfg, mod.Slug)
		if err != nil {
			logger.Log.Printf(
				"Error getting versions for mod %s: %v",
				mod.Name,
				err,
			)
			return err
		}

		mu.Lock()
		mods = append(mods, ModSearch{
			ID:          mod.Slug,
			Name:        mod.Name,
			Description: mod.Summary,
			Downloads:   formatDownloads(mod.DownloadCount),
			URL:         mod.Links.WebsiteURL,
			ClientSide:  "",
			ServerSide:  "",
			Versions:    versions,
		})
		mu.Unlock()

		return nil
	}

	jobs := make(chan curseforgeMod, len(response.Data))

	results := util.WorkerPool(
		jobs,
		processMod,
		len(response.Data),
	)

	go func() {
		for _, mod := range response.Data {
			jobs <- mod
		}

		close(jobs)
	}()

	for err := range results {
		if err != nil {
			logger.Log.Printf("Error processing mod: %v", err)
			return nil, err
		}
	}

	return mods, nil
}

func getDownloadURLCurseforge(cfg *config.Config, id, version string) (string, error) {
	logger.Log.Printf("Getting download URL for mod: %s", id)

	mod, err := getModBySlug(id)
	if err != nil {
		return "", err
	}

	files, err := getModFiles(cfg, mod.ID)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if file.DisplayName == version {
			if file.DownloadURL == "" {
				return "", fmt.Errorf("download url is empty")
			}

			logger.Log.Printf(
				"Download URL obtained for %s: %s",
				version,
				file.DownloadURL,
			)

			return file.DownloadURL, nil
		}
	}

	return "", fmt.Errorf("version not found: %s", version)
}

func getModVersionsCurseforge(cfg *config.Config, id string) ([]string, error) {
	mod, err := getModBySlug(id)
	if err != nil {
		return nil, err
	}

	files, err := getModFiles(cfg, mod.ID)
	if err != nil {
		return nil, err
	}

	var versions []string

	for _, file := range files {
		versions = append(versions, file.DisplayName)
	}

	if len(versions) == 0 {
		return nil, fmt.Errorf("no versions found")
	}

	return versions, nil
}

func getLatestVersionCurseforge(cfg *config.Config, id string) (string, error) {
	versions, err := getModVersionsCurseforge(cfg, id)
	if err != nil {
		return "", err
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("no versions found")
	}

	return versions[0], nil
}

func getModBySlug(slug string) (*curseforgeMod, error) {
	params := url.Values{}
	params.Set("gameId", strconv.Itoa(minecraftGame))
	params.Set("slug", slug)

	var response curseforgeSearchResponse

	err := makeCurseforgeRequest(
		"/mods/search?"+params.Encode(),
		&response,
	)
	if err != nil {
		return nil, err
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("mod not found")
	}

	return &response.Data[0], nil
}

func getModFiles(cfg *config.Config, modID int) ([]curseforgeFile, error) {
	params := url.Values{}
	params.Set("pageSize", "50")

	if cfg.Minecraft != "" {
		params.Set("gameVersion", cfg.Minecraft)
	}

	modLoaderType := curseforgeLoaderType(cfg.Loader)
	if modLoaderType != 0 {
		params.Set("modLoaderType", strconv.Itoa(modLoaderType))
	}

	var response curseforgeFilesResponse

	err := makeCurseforgeRequest(
		fmt.Sprintf(
			"/mods/%d/files?%s",
			modID,
			params.Encode(),
		),
		&response,
	)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

func formatDownloads(downloads float64) string {
	switch {
	case downloads >= 1_000_000:
		return fmt.Sprintf("%.1fM", downloads/1_000_000)

	case downloads >= 1_000:
		return fmt.Sprintf("%.1fK", downloads/1_000)

	default:
		return fmt.Sprintf("%.0f", downloads)
	}
}
