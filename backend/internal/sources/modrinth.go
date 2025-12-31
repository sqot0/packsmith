package sources

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"sync"

	"github.com/sqot0/packsmith/backend/internal/config"
	"github.com/sqot0/packsmith/backend/internal/logger"
	"github.com/sqot0/packsmith/backend/internal/util"
)

type ModrinthSearchMod struct {
	Slug, Title, Description string
	ClientSide               string `json:"client_side"`
	ServerSide               string `json:"server_side"`
	Downloads                int
}

type ModrinthSearch struct {
	Hits []ModrinthSearchMod
}

type ModrinthModVersion struct {
	GameVersions []string `json:"game_versions"`
	Loaders      []string `json:"loaders"`
	Version      string   `json:"version_number"`
	Files        []struct {
		URL     string `json:"url"`
		Primary bool   `json:"primary"`
	}
}

func searchModsModrinth(cfg *config.Config, query string) ([]ModSearch, error) {
	logger.Log.Printf("Searching Modrinth for query: %s", query)
	params := url.Values{
		"query":  {query},
		"limit":  {"5"},
		"facets": {fmt.Sprintf(`[["project_type:mod"],["versions:%s"],["categories:%s"]]`, cfg.Minecraft, cfg.Loader)},
	}

	req, _ := http.NewRequest("GET", "https://api.modrinth.com/v2/search?"+params.Encode(), nil)
	setHeadersForRequest(req)

	logger.Log.Printf("Making HTTP request to Modrinth search API")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log.Printf("Error making request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var data ModrinthSearch

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		logger.Log.Printf("Error decoding JSON response: %v", err)
		return nil, err
	}

	var mu sync.Mutex
	mods := make([]ModSearch, 0, len(data.Hits))

	processMod := func(h ModrinthSearchMod) error {
		versionsUrl := fmt.Sprintf("https://api.modrinth.com/v2/project/%s/version", h.Slug)
		req, _ := http.NewRequest("GET", versionsUrl, nil)
		setHeadersForRequest(req)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			logger.Log.Printf("Error making request for versions: %v", err)
			return err
		}
		defer resp.Body.Close()

		var versionsData []ModrinthModVersion

		if err := json.NewDecoder(resp.Body).Decode(&versionsData); err != nil {
			logger.Log.Printf("Error decoding JSON response for versions: %v", err)
			return err
		}

		var versions []string
		for _, v := range versionsData {
			if slices.Contains(v.GameVersions, cfg.Minecraft) &&
				slices.Contains(v.Loaders, cfg.Loader) {
				versions = append(versions, v.Version)
			}
		}

		mu.Lock()
		mods = append(mods, ModSearch{
			ID: h.Slug, Name: h.Title, Description: h.Description,
			ClientSide: h.ClientSide, ServerSide: h.ServerSide,
			Downloads: strconv.Itoa(h.Downloads),
			URL:       "https://modrinth.com/mod/" + h.Slug,
			Versions:  versions,
		})
		mu.Unlock()
		return nil
	}

	jobs := make(chan ModrinthSearchMod, len(data.Hits))
	results := util.WorkerPool(jobs, processMod, len(data.Hits))

	go func() {
		for _, hit := range data.Hits {
			jobs <- hit
		}
		close(jobs)
	}()

	for err := range results {
		if err != nil {
			logger.Log.Printf("Error processing mod: %v", err)
			return nil, err
		}
	}

	logger.Log.Printf("Found %d mods on Modrinth", len(mods))
	return mods, nil
}

func getDownloadURLModrinth(cfg *config.Config, id, version string) (string, error) {
	logger.Log.Printf("Getting download URL for Modrinth mod: %s", id)
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("https://api.modrinth.com/v2/project/%s/version", id), nil)
	setHeadersForRequest(req)

	logger.Log.Printf("Making HTTP request to Modrinth version API")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log.Printf("Error making request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	var versions []ModrinthModVersion

	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		logger.Log.Printf("Error decoding JSON response: %v", err)
		return "", err
	}

	for _, v := range versions {
		if v.Version == version && slices.Contains(v.GameVersions, cfg.Minecraft) &&
			slices.Contains(v.Loaders, cfg.Loader) {
			for _, f := range v.Files {
				if f.Primary {
					logger.Log.Printf("Download URL obtained: %s, version: %s", f.URL, v.Version)
					return f.URL, nil
				}
			}
		}
	}
	logger.Log.Println("No compatible version found")
	return "", fmt.Errorf("no compatible version found")
}

func getModVersionsModrinth(cfg *config.Config, id string) ([]string, error) {
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("https://api.modrinth.com/v2/project/%s/version", id), nil)
	setHeadersForRequest(req)

	logger.Log.Printf("Making HTTP request to Modrinth version API")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log.Printf("Error making request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var versions []ModrinthModVersion

	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		logger.Log.Printf("Error decoding JSON response: %v", err)
		return nil, err
	}
	var compatibleVersions []string
	for _, v := range versions {
		if slices.Contains(v.GameVersions, cfg.Minecraft) &&
			slices.Contains(v.Loaders, cfg.Loader) {
			compatibleVersions = append(compatibleVersions, v.Version)
		}
	}
	logger.Log.Printf("Found %d compatible versions", len(compatibleVersions))
	return compatibleVersions, nil
}

func getLatestVersionModrinth(cfg *config.Config, id string) (string, error) {
	logger.Log.Printf("Getting latest version for Modrinth mod: %s", id)
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("https://api.modrinth.com/v2/project/%s/version", id), nil)
	setHeadersForRequest(req)

	logger.Log.Printf("Making HTTP request to Modrinth version API")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log.Printf("Error making request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	var versions []ModrinthModVersion

	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		logger.Log.Printf("Error decoding JSON response: %v", err)
		return "", err
	}

	for _, v := range versions {
		if slices.Contains(v.GameVersions, cfg.Minecraft) &&
			slices.Contains(v.Loaders, cfg.Loader) {
			logger.Log.Printf("Latest version found: %s", v.Version)
			return v.Version, nil
		}
	}
	logger.Log.Println("No compatible version found")
	return "", fmt.Errorf("no compatible version found")
}
