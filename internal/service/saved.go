package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend_api/internal/model"
	"github.com/VladPetriv/scanner_backend_api/internal/store"
)

type SavedDBService struct {
	store *store.Store
}

func NewSavedService(store *store.Store) *SavedDBService {
	return &SavedDBService{store: store}
}

func (s *SavedDBService) GetSavedMessages(ID int) ([]model.Saved, error) {
	savedMessages, err := s.store.Saved.GetSavedMessages(ID)
	if err != nil {
		return nil, fmt.Errorf("[Saved] srv.GetSavedMessages error: %w", err)
	}

	return savedMessages, nil
}
func (s *SavedDBService) CreateSavedMessage(savedMessage *model.Saved) error {
	_, err := s.store.Saved.CreateSavedMessage(savedMessage)
	if err != nil {
		return fmt.Errorf("[Saved] srv.CreateSavedMessage error: %w", err)
	}

	return nil
}

func (s *SavedDBService) DeleteSavedMessage(ID int) error {
	_, err := s.store.Saved.DeleteSavedMessage(ID)
	if err != nil {
		return fmt.Errorf("[Saved] srv.DeleteSavedMessage error: %w", err)
	}

	return nil
}
