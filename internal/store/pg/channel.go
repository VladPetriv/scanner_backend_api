package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend_api/internal/model"
)

var (
	ErrChannelsCountNotFound = errors.New("channels count not found")
	ErrChannelsNotFound      = errors.New("channels not found")
	ErrChannelNotFound       = errors.New("channel not found")
)

type ChannelRepo struct {
	db *DB
}

func NewChannelRepo(db *DB) *ChannelRepo {
	return &ChannelRepo{db: db}
}

func (c *ChannelRepo) CreateChannel(channel *model.ChannelDTO) error {
	_, err := c.db.Exec(
		"INSERT INTO channel(name, title, imageurl) VALUES ($1, $2, $3);",
		channel.Name, channel.Title, channel.ImageURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}

	return nil
}

func (c *ChannelRepo) GetChannelsCount() (int, error) {
	var count int

	err := c.db.Get(&count, "SELECT COUNT(*) FROM channel;")
	if err == sql.ErrNoRows {
		return 0, ErrChannelsCountNotFound
	}

	if err != nil {
		return count, fmt.Errorf("failed to get count of channels: %w", err)
	}

	return count, nil
}

func (c *ChannelRepo) GetChannelsByPage(offset int) ([]model.Channel, error) {
	channels := make([]model.Channel, 0, 10)

	err := c.db.Select(&channels, "SELECT * FROM channel OFFSET $1 LIMIT 10;", offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get channels by page: %w", err)
	}

	if len(channels) == 0 {
		return nil, ErrChannelsNotFound
	}

	return channels, nil
}

func (c *ChannelRepo) GetChannelByName(name string) (*model.Channel, error) {
	var channel model.Channel

	err := c.db.Get(&channel, "SELECT * FROM channel WHERE name = $1;", name)
	if err == sql.ErrNoRows {
		return nil, ErrChannelNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get channel by name: %w", err)
	}

	return &channel, nil
}
