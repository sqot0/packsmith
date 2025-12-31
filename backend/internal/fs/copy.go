package fs

import (
	"io"
	"os"

	"github.com/sqot0/packsmith/backend/internal/logger"
)

func Copy(src, dst string) error {
	logger.Log.Printf("Copying file from %s to %s", src, dst)
	in, err := os.Open(src)
	if err != nil {
		logger.Log.Printf("Error opening source file: %v", err)
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		logger.Log.Printf("Error creating destination file: %v", err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		logger.Log.Printf("Error copying file content: %v", err)
		return err
	}

	err = out.Sync()
	if err != nil {
		logger.Log.Printf("Error syncing file: %v", err)
		return err
	}
	logger.Log.Println("File copied successfully")
	return nil
}
