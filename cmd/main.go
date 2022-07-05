package main

import (
	"go.uber.org/zap"

	"github.com/VladPetriv/scanner_backend_api/internal/handler"
	"github.com/VladPetriv/scanner_backend_api/internal/server"
	"github.com/VladPetriv/scanner_backend_api/internal/service"
	"github.com/VladPetriv/scanner_backend_api/internal/store"
	"github.com/VladPetriv/scanner_backend_api/pkg/config"
	"github.com/VladPetriv/scanner_backend_api/pkg/logger"
)

func main() {
	log := logger.Get()

	cfg, err := config.Get()
	if err != nil {
		log.Error("failed to load config", zap.Error(err))
	}

	store, err := store.New(cfg, log)
	if err != nil {
		log.Error("failed to create store", zap.Error(err))
	}

	service, err := service.New(store, cfg.JwtSecretKey)
	if err != nil {
		log.Error("failed to create service", zap.Error(err))
	}

	handler := handler.New(service, log)

	server := new(server.Server)

	log.Info("start server", zap.String("PORT", cfg.Port))

	if err := server.Start(handler.InitRoutes()); err != nil {
		log.Error("failed to start server", zap.Error(err))
	}
}
