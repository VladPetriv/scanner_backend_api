package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend_api/internal/model"
	"github.com/VladPetriv/scanner_backend_api/internal/store"
	"github.com/VladPetriv/scanner_backend_api/pkg/utils"
)

type MessageDBService struct {
	store *store.Store
}

func NewMessageService(store *store.Store) *MessageDBService {
	return &MessageDBService{store: store}
}

func (m *MessageDBService) GetMessagesCount() (int, error) {
	count, err := m.store.Message.GetMessagesCount()
	if err != nil {
		return 0, fmt.Errorf("[Message] srv.GetMessagesCount error: %w", err)
	}

	return count, nil
}

func (m *MessageDBService) GetMessagesCountByChannelID(ID int) (int, error) {
	count, err := m.store.Message.GetMessagesCountByChannelID(ID)
	if err != nil {
		return 0, fmt.Errorf("[Message] srv.GetMessagesCountByChannelID error: %w", err)
	}

	return count, nil
}

func (m *MessageDBService) GetFullMessagesByPage(page int) ([]model.FullMessage, error) {
	messages, err := m.store.Message.GetFullMessagesByPage(utils.FormatPage(page))
	if err != nil {
		return nil, fmt.Errorf("[Message] srv.GetFullMessagesByPage error: %w", err)
	}

	return messages, nil
}

func (m *MessageDBService) GetFullMessagesByChannelIDAndPage(ID, page int) ([]model.FullMessage, error) {
	messages, err := m.store.Message.GetFullMessagesByChannelIDAndPage(ID, utils.FormatPage(page))
	if err != nil {
		return nil, fmt.Errorf("[Message] srv.GetFullMessagesByChannelIDAndPage error: %w", err)
	}

	return messages, nil
}

func (m *MessageDBService) GetFullMessagesByUserID(ID int) ([]model.FullMessage, error) {
	messages, err := m.store.Message.GetFullMessagesByUserID(ID)
	if err != nil {
		return nil, fmt.Errorf("[Message] srv.GetFullMessagesByUserID error: %w", err)
	}

	return messages, nil
}

func (m *MessageDBService) GetFullMessageByID(ID int) (*model.FullMessage, error) {
	message, err := m.store.Message.GetFullMessageByID(ID)
	if err != nil {
		return nil, fmt.Errorf("[Message] srv.GetFullMessageByID error: %w", err)
	}

	return message, nil
}
