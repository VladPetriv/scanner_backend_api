package service

import (
	"github.com/VladPetriv/scanner_backend_api/internal/store"
)

type UserDBService struct {
	store *store.Store
}

func NewUserService(store *store.Store) *UserDBService {
	return &UserDBService{store: store}
}
