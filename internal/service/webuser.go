package service

import (
	"github.com/VladPetriv/scanner_backend_api/internal/store"
)

type WebUserDBService struct {
	store *store.Store
}

func NewWebUserService(store *store.Store) *WebUserDBService {
	return &WebUserDBService{store: store}
}
