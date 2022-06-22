package service

import (
	"github.com/VladPetriv/scanner_backend_api/internal/store"
)

type MessageDBService struct {
	store *store.Store
}

func NewMessageService(store *store.Store) *MessageDBService {
	return &MessageDBService{store: store}
}
