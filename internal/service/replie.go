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

func (r *ReplieDBService) CreateReplie(replie *model.ReplieDTO) error {
	err := r.store.Replie.CreateReplie(replie)
	if err != nil {
		return fmt.Errorf("[Replie] srv.CreateReplie error: %w", err)
	}

	return nil
}

func (r *ReplieDBService) GetFullRepliesByMessageID(ID int) ([]model.FullReplie, error) {
	replies, err := r.store.Replie.GetFullRepliesByMessageID(ID)
	if err != nil {
		return nil, fmt.Errorf("[Replie] srv.GetFullRepliesByMessageID error: %w", err)
	}

	return replies, nil
}
