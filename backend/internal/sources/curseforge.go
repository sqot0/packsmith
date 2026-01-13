package sources

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/sqot0/packsmith/backend/internal/config"
	"github.com/sqot0/packsmith/backend/internal/logger"
	"github.com/sqot0/packsmith/backend/internal/util"
)

func searchModsCurseforge(cfg *config.Config, query string) ([]ModSearch, error) {
	logger.Log.Printf("Searching CurseForge for query: %s", query)
	params := url.Values{}
	params.Set("page", "1")
	params.Set("pageSize", "5")
	params.Set("sortBy", "relevancy")
	params.Set("class", "mc-mods")

	params.Set("version", cfg.Minecraft)
	params.Set("search", query)
	switch cfg.Loader {
	case "forge":
		params.Set("gameVersionTypeId", "1")
	case "fabric":
		params.Set("gameVersionTypeId", "4")
	case "neoforge":
		params.Set("gameVersionTypeId", "6")
	case "quilt":
		params.Set("gameVersionTypeId", "5")
	}

	req, _ := http.NewRequest("GET", "https://www.curseforge.com/minecraft/search?"+params.Encode(), nil)
	setHeadersForRequest(req)

	logger.Log.Printf("Making HTTP request to CurseForge search")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log.Printf("Error making request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		logger.Log.Printf("CurseForge search returned status %d", resp.StatusCode)
		return nil, fmt.Errorf("curseforge search returned status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logger.Log.Printf("Error parsing HTML: %v", err)
		return nil, err
	}

	type curseforgeModInfo struct {
		id, name, description, downloads, url string
	}

	var modInfos []curseforgeModInfo
	doc.Find(".project-card").Each(func(i int, s *goquery.Selection) {
		modUrl := "https://www.curseforge.com" + s.Find(".name").AttrOr("href", "")
		if modUrl == "" {
			return
		}
		id := strings.TrimPrefix(modUrl, "https://www.curseforge.com/minecraft/mc-mods/")
		name := s.Find(".name span").Text()
		description := s.Find(".description").Text()
		downloads := s.Find(".details-list .detail-downloads").Text()

		modInfos = append(modInfos, curseforgeModInfo{
			id: id, name: name, description: description, downloads: downloads, url: modUrl,
		})
	})

	var mods []ModSearch
	var mu sync.Mutex
	processMod := func(info curseforgeModInfo) error {
		versions, err := getModVersionsCurseforge(cfg, info.id)
		if err != nil {
			logger.Log.Printf("Error getting version for mod %s: %v", info.id, err)
			return err
		}
		mu.Lock()
		mods = append(mods, ModSearch{
			ID:          info.id,
			Name:        info.name,
			Description: info.description,
			Downloads:   info.downloads,
			URL:         info.url,
			ClientSide:  "",
			ServerSide:  "",
			Versions:    versions,
		})
		mu.Unlock()
		return nil
	}

	jobs := make(chan curseforgeModInfo, len(modInfos))
	results := util.WorkerPool(jobs, processMod, len(modInfos))

	go func() {
		for _, info := range modInfos {
			jobs <- info
		}
		close(jobs)
	}()

	for err := range results {
		if err != nil {
			return nil, err
		}
	}

	logger.Log.Printf("Found %d mods on CurseForge", len(mods))
	return mods, nil
}

func getDownloadURLCurseforge(cfg *config.Config, id, version string) (string, error) {
	logger.Log.Printf("Getting download URL for CurseForge mod: %s", id)
	params := url.Values{}
	params.Set("page", "1")
	params.Set("pageSize", "20")
	params.Set("showAlphaFiles", "hide")
	params.Set("class", "mc-mods")

	params.Set("version", cfg.Minecraft)
	switch cfg.Loader {
	case "forge":
		params.Set("gameVersionTypeId", "1")
	case "fabric":
		params.Set("gameVersionTypeId", "4")
	case "neoforge":
		params.Set("gameVersionTypeId", "6")
	case "quilt":
		params.Set("gameVersionTypeId", "5")
	}

	searchUrl := fmt.Sprintf("https://www.curseforge.com/minecraft/mc-mods/%s/files/all?", id)
	req, _ := http.NewRequest("GET", searchUrl+params.Encode(), nil)
	setHeadersForRequest(req)

	logger.Log.Printf("Making HTTP request to CurseForge files page")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log.Printf("Error making request: %v", err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		logger.Log.Printf("CurseForge files page returned status %d", resp.StatusCode)
		return "", fmt.Errorf("curseforge search returned status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logger.Log.Printf("Error parsing HTML: %v", err)
		return "", err
	}

	projectId := doc.Find(".project-id").Text()
	if projectId == "" {
		logger.Log.Println("Could not find project ID")
		return "", fmt.Errorf("could not find project ID")
	}

	var downloadUrl string
	found := false
	doc.Find(".file-row-details").Each(func(i int, s *goquery.Selection) {
		title := s.Find("span.name").AttrOr("title", "")
		if title == version {
			fileDetailsURL := s.AttrOr("href", "")
			if fileDetailsURL != "" {
				fileId := fileDetailsURL[strings.LastIndex(fileDetailsURL, "/")+1:]
				downloadUrl = fmt.Sprintf("https://www.curseforge.com/api/v1/mods/%s/files/%s/download", projectId, fileId)
				found = true
				return
			}
		}
	})

	if !found {
		logger.Log.Printf("Could not find file for version: %s", version)
		return "", fmt.Errorf("could not find file for version: %s", version)
	}

	logger.Log.Printf("Download URL obtained: %s, version: %s", downloadUrl)
	return downloadUrl, nil
}

func getModVersionsCurseforge(cfg *config.Config, id string) ([]string, error) {
	params := url.Values{}
	params.Set("page", "1")
	params.Set("pageSize", "20")
	params.Set("showAlphaFiles", "hide")
	params.Set("class", "mc-mods")

	params.Set("version", cfg.Minecraft)
	switch cfg.Loader {
	case "forge":
		params.Set("gameVersionTypeId", "1")
	case "fabric":
		params.Set("gameVersionTypeId", "4")
	case "neoforge":
		params.Set("gameVersionTypeId", "6")
	case "quilt":
		params.Set("gameVersionTypeId", "5")
	}

	searchUrl := fmt.Sprintf("https://www.curseforge.com/minecraft/mc-mods/%s/files/all?", id)
	req, _ := http.NewRequest("GET", searchUrl+params.Encode(), nil)
	setHeadersForRequest(req)

	logger.Log.Printf("Making HTTP request to CurseForge files page for version")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log.Printf("Error making request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		logger.Log.Printf("CurseForge files page returned status %d", resp.StatusCode)
		return nil, fmt.Errorf("curseforge files page returned status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logger.Log.Printf("Error parsing HTML: %v", err)
		return nil, err
	}

	fileVersions := doc.Find("span.name")
	if fileVersions.Length() == 0 {
		logger.Log.Println("Could not find file version")
		return nil, fmt.Errorf("could not find file version")
	}

	var result []string
	fileVersions.Each(func(i int, s *goquery.Selection) {
		version := s.AttrOr("title", "")
		if version != "" {
			result = append(result, version)
		}
	})

	return result, nil
}

func getLatestVersionCurseforge(cfg *config.Config, id string) (string, error) {
	versions, err := getModVersionsCurseforge(cfg, id)
	if err != nil {
		return "", err
	}
	if len(versions) == 0 {
		logger.Log.Println("No versions found for mod")
		return "", fmt.Errorf("no versions found for mod")
	}

	return versions[0], nil
}
