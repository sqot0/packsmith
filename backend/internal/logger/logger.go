package logger

import (
	"log"
	"os"
	"path"
)

var Log *log.Logger

func Init() {
	Log = log.New(os.Stdout, "[APP] ", log.Ldate|log.Ltime)
	Log.Println("Initializing logger")

	logFile := getLogFilePath()
	Log.Printf("Log file path: %s", logFile)

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		Log.Printf("Error opening log file: %v", err)
		panic(err)
	}

	Log = log.New(file, "[APP] ", log.Ldate|log.Ltime)
	Log.Println("Logger initialized successfully")
}

func getLogFilePath() string {
	cwd, err := os.Getwd()
	if err != nil {
		Log.Printf("Error getting current working directory: %v", err)
		panic(err)
	}
	return path.Join(cwd, "app.log")
}
