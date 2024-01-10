package log_config

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var defaultLogFilePath string

func createLogFolderIfNotExist(path string) {

	dir := filepath.Dir(path)

	if _, errDir := os.Stat(dir); os.IsNotExist(errDir) {

		log.Println("Creating", dir)

		errMkdir := os.MkdirAll(dir, 0664)

		if errMkdir != nil {
			log.Println("Fail to create directory", dir)
		} else {
			log.Println("Success to create directory", dir)

		}
	}
}

func openOrCreateLogFile(path string) (*os.File, error) {

	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {

		logFile, errCreateFile := os.Create(path)

		if errCreateFile != nil {
			log.Println("Cant create log file", errCreateFile)
		}

		return logFile, nil

	}

	return logFile, nil

}

func DefaultLogging(path ...string) {

	gin.DisableConsoleColor()

	if len(path) > 0 && path[0] != "" {
		defaultLogFilePath = path[0]
	} else {
		defaultLogFilePath = "logs/gin.log"
	}

	createLogFolderIfNotExist(defaultLogFilePath)

	f, _ := openOrCreateLogFile(defaultLogFilePath)

	gin.DefaultWriter = io.MultiWriter(f)

	log.SetOutput(gin.DefaultWriter)

}
