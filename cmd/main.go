package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	logger := setupLogger()

	logger.SetReportCaller(true)
}

func setupLogger() *log.Logger {
	logger := log.New()
	logger.SetReportCaller(true)
	return logger
}
