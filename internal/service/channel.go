package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend_api/internal/model"
	"github.com/VladPetriv/scanner_backend_api/internal/store"
	"github.com/VladPetriv/scanner_backend_api/pkg/utils"
)

type ChannelDBService struct {
	store *store.Store
}

func NewChannelService(store *store.Store) *ChannelDBService {
	return &ChannelDBService{store: store}
}

func (c *ChannelDBService) CreateChannel(channel *model.ChannelDTO) error {
	err := c.store.Channel.CreateChannel(channel)
	if err != nil {
		return fmt.Errorf("[Channel] srv.CreateChannel error: %w", err)
	}

	return nil
}

func (c *ChannelDBService) GetChannelsCount() (int, error) {
	count, err := c.store.Channel.GetChannelsCount()
	if err != nil {
		return 0, fmt.Errorf("[Channel] srv.GetChannelsCount error: %w", err)
	}

	return count, nil
}

func (c *ChannelDBService) GetChannelsByPage(page int) ([]model.Channel, error) {
	channels, err := c.store.Channel.GetChannelsByPage(utils.FormatPage(page))
	if err != nil {
		return nil, fmt.Errorf("[Channel] srv.GetChannelsByPage error: %w", err)
	}

	return channels, nil
}

func (c *ChannelDBService) GetChannelByName(name string) (*model.Channel, error) {
	channel, err := c.store.Channel.GetChannelByName(name)
	if err != nil {
		return nil, fmt.Errorf("[Channel] srv.GetChannelByName error: %w", err)
	}

	return channel, nil
}
