package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend_api/internal/model"
	"github.com/VladPetriv/scanner_backend_api/internal/store"
)

type ReplieDBService struct {
	store *store.Store
}

func NewReplieService(store *store.Store) *ReplieDBService {
	return &ReplieDBService{store: store}
}

func (r *ReplieDBService) GetFullRepliesByMessageID(ID int) ([]model.FullReplie, error) {
	replies, err := r.store.Replie.GetFullRepliesByMessageID(ID)
	if err != nil {
		return nil, fmt.Errorf("[Replie] srv.GetFullRepliesByMessageID error: %w", err)
	}

	return replies, nil
}
