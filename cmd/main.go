package main

import (
	"go.uber.org/zap"

	"github.com/VladPetriv/scanner_backend_api/internal/handler"
	"github.com/VladPetriv/scanner_backend_api/internal/server"
	"github.com/VladPetriv/scanner_backend_api/internal/service"
	"github.com/VladPetriv/scanner_backend_api/internal/store"
	"github.com/VladPetriv/scanner_backend_api/internal/store/kafka"
	"github.com/VladPetriv/scanner_backend_api/pkg/config"
	"github.com/VladPetriv/scanner_backend_api/pkg/logger"
)

// @title        Scanner Back-End API
// @version      1.0
// @description  Back-End side for telegram scanner
// @produce      json

// @host      localhost:3000
// @BasePath  /

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
func main() {
	cfg, err := config.Get()
	if err != nil {
		panic(err)
	}

	log := logger.Get(cfg.LogLevel)

	store, err := store.New(cfg, log)
	if err != nil {
		log.Error("failed to create store", zap.Error(err))
	}

	service, err := service.New(store, cfg.JwtSecretKey)
	if err != nil {
		log.Error("failed to create service", zap.Error(err))
	}

	go kafka.GetChannelFromQueue(service, cfg, log)
	go kafka.GetDataFromQueue(service, cfg, log)

	server := new(server.Server)

	handler := handler.New(service, log)

	log.Info("start server", zap.String("PORT", cfg.Port))

	if err := server.Start(handler.InitRoutes()); err != nil {
		log.Error("failed to start server", zap.Error(err))
	}
}
