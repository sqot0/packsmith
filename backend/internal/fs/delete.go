package fs

import (
	"os"
	"path/filepath"

	"github.com/sqot0/packsmith/backend/internal/logger"
)

func Delete(projectPath, filename string) error {
	logger.Log.Printf("Deleting file: %s in project: %s", filename, projectPath)
	cacheModPath := filepath.Join(projectPath, cacheDir, filename)
	err := os.RemoveAll(cacheModPath)
	if err != nil {
		logger.Log.Printf("Error deleting file: %v", err)
		return err
	}
	logger.Log.Println("File deleted successfully")
	return nil
}
