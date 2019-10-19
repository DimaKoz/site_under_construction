package main

import (
	"github.com/google/logger"
	"os"
	"under_construction/app"
)

var logFile *os.File
var log *logger.Logger

func initLogger() {
	var errLog error
	logFile, errLog = os.OpenFile(app.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if errLog != nil {
		logger.Fatalf("Failed to open log file: %v", errLog)
	}
	log = logger.Init("LoggerExample", true, false, logFile)
	logger.Warningln("")
	logger.Warningln("================================================================")
	logger.Warningln("")
	logger.Warningln("Logger started")
	logger.Warningln("")
	logger.Warningln("================================================================")

}

func closeLogger() {
	if log != nil {
		log.Close()
	}

	if logFile != nil {
		_ = logFile.Close()
	}
}
