package fs

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/sqot0/packsmith/backend/internal/logger"
)

const cacheDir = "cache"

func Download(projectPath, fileURL, version string) (string, error) {
	logger.Log.Printf("Downloading file from URL: %s", fileURL)
	resp, err := http.Get(fileURL)
	if err != nil {
		logger.Log.Printf("Error making HTTP request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Log.Printf("Download failed with status: %s", resp.Status)
		return "", fmt.Errorf("download failed: %s", resp.Status)
	}

	cacheFolder := filepath.Join(projectPath, cacheDir)
	logger.Log.Printf("Ensuring cache folder exists: %s", cacheFolder)
	if err := os.MkdirAll(cacheFolder, 0o755); err != nil {
		logger.Log.Printf("Error creating cache folder: %v", err)
		return "", err
	}

	name := getFilename(resp, version)
	if name == "" {
		logger.Log.Println("Could not determine filename")
		return "", fmt.Errorf("could not determine filename")
	}
	logger.Log.Printf("Determined filename: %s", name)

	out, err := os.Create(filepath.Join(cacheFolder, name))
	if err != nil {
		logger.Log.Printf("Error creating file: %v", err)
		return "", err
	}
	defer out.Close()

	logger.Log.Println("Copying file content")
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logger.Log.Printf("Error copying file content: %v", err)
		return "", err
	}
	logger.Log.Println("File downloaded successfully")
	return name, nil
}

func getFilename(resp *http.Response, version string) string {
	logger.Log.Println("Determining filename from response")
	fileURL := resp.Request.URL
	if filepath.Base(fileURL.Path) != "download" {
		logger.Log.Printf("Filename from URL: %s", filepath.Base(fileURL.Path))
		return filepath.Base(fileURL.Path)
	}

	if version != "" && path.Ext(version) == ".jar" {
		logger.Log.Printf("Filename from version: %s", version)
		return version
	}

	logger.Log.Println("Could not determine filename")
	return ""
}
