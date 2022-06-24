package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend_api/internal/model"
	"github.com/VladPetriv/scanner_backend_api/internal/store"
)

type ChannelDBService struct {
	store *store.Store
}

func NewChannelService(store *store.Store) *ChannelDBService {
	return &ChannelDBService{store: store}
}

func (c *ChannelDBService) GetChannelsCount() (int, error) {
	count, err := c.store.Channel.GetChannelsCount()
	if err != nil {
		return 0, fmt.Errorf("[Channel] srv.GetChannelsCount error: %w", err)
	}

	return count, nil
}

func (c *ChannelDBService) GetChannelsByPage(page int) ([]model.Channel, error) {
	if page == 0 || page == 1 {
		page = 0
	} else if page == 2 {
		page = 10
	} else {
		page *= 10
		page -= 10
	}

	channels, err := c.store.Channel.GetChannelsByPage(page)
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
