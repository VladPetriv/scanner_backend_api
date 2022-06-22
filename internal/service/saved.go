package service

import (
	"github.com/VladPetriv/scanner_backend_api/internal/store"
)

type SavedDBService struct {
	store *store.Store
}

func NewSavedService(store *store.Store) *SavedDBService {
	return &SavedDBService{store: store}
}
