package main

import (
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/pedroxer/auth-service/internal/config"
	"github.com/pedroxer/auth-service/internal/database"
	"github.com/pedroxer/auth-service/internal/routes"
	"github.com/pedroxer/auth-service/internal/service/auth"
	storage2 "github.com/pedroxer/auth-service/internal/storage"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := setupLogger()
	data, err := os.ReadFile("../config/config.json")
	if err != nil {
		logger.Fatal(err)
	}

	cfg := new(config.Config)
	if err := json.Unmarshal(data, cfg); err != nil {
		logger.Fatal(err)
	}
	logger.Infof("Success read config.")
	if err := env.Parse(cfg); err != nil {
		logger.Fatal(err)
	}
	logger.Infof("Success read envs.")

	psDb, err := database.ConnectToPG(&cfg.Postgres)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("connected to db '%s'", cfg.Postgres.Db)

	storage := storage2.NewStorage(logger, psDb)

	authService := auth.NewUserAuth(storage, storage, logger)
	logger.Info("created auth service")
	if storage == nil {
		logger.Fatal("storage is nil")
	}
	if authService == nil {
		logger.Fatal("auth service is nil")
	}
	if psDb == nil {
		logger.Fatal("db is nil")
	}
	router := routes.New(cfg, logger, authService)
	fmt.Println(router)
	errChan := make(chan error, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	// go func() {
	// 	if err := router.Start(); err != nil {
	// 		errChan <- err
	// 	}
	// }()
		if err := router.Start(); err != nil {
			logger.Fatal(err)
		}

	var startErr error
	select {
	case sig := <-sigChan:
		logger.Infoln("got signal:", sig)
	case err := <-errChan:
		logger.Info(err)
		startErr = err
	}

	logger.Info("gracefully shutting down the server...")
	if err := router.Shutdown(); err != nil {
		logger.Warn(err)
	}

	logger.Info("closing the databases...")
	if err := psDb.Close(); err != nil {
		logger.Warn(err)
	}

	if startErr != nil {
		logger.Fatal(startErr)
	}

	logger.Info("exited")

}

func setupLogger() *log.Logger {
	logger := log.New()
	logger.SetReportCaller(true)
	return logger
}
