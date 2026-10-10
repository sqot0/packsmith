package logger

import (
	"log"
	"os"
	"path"
)

var Log = log.New(os.Stderr, "[APP] ", log.Ldate|log.Ltime)

func Init() {
	logFile := getLogFilePath()
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		Log = log.New(os.Stderr, "[APP] ", log.Ldate|log.Ltime)
		Log.Printf("Error opening log file: %v", err)
		return
	}

	Log = log.New(file, "[APP] ", log.Ldate|log.Ltime)
	Log.Println("Logger initialized successfully")
}

func getLogFilePath() string {
	cwd, err := os.Getwd()
	if err != nil {
		Log.Printf("Error getting current working directory: %v", err)
		return "app.log"
	}
	return path.Join(cwd, "app.log")
}
