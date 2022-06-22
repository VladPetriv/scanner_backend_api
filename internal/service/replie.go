package service

import (
	"github.com/VladPetriv/scanner_backend_api/internal/store"
)

type ReplieDBService struct {
	store *store.Store
}

func NewReplieService(store *store.Store) *ReplieDBService {
	return &ReplieDBService{store: store}
}
