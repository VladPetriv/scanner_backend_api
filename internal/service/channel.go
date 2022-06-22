package service

import (
	"github.com/VladPetriv/scanner_backend_api/internal/store"
)

type ChannelDBService struct {
	store *store.Store
}

func NewChannelService(store *store.Store) *ChannelDBService {
	return &ChannelDBService{store: store}
}
