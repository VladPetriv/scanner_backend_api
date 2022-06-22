package main

import (
	"github.com/VladPetriv/scanner_backend_api/internal/service"
	"github.com/VladPetriv/scanner_backend_api/internal/store"
	"github.com/VladPetriv/scanner_backend_api/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	log := logger.Get()

	store, err := store.New(log)
	if err != nil {
		log.Error("failed to create store", zap.Error(err))
	}

	service, err := service.New(store)
	if err != nil {
		log.Error("failed to create service", zap.Error(err))
	}
}
