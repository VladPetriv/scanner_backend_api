package store

import "github.com/VladPetriv/scanner_backend_api/internal/model"

//go:generate mockery --dir . --name ChannelRepo --output ./mocks
type ChannelRepo interface {
	GetChannelsCount() (int, error)
	GetChannelsByPage(offset int) ([]model.Channel, error)
	GetChannelByName(name string) (*model.Channel, error)
}

//go:generate mockery --dir . --name UserRepo --output ./mocks
type UserRepo interface{}

//go:generate mockery --dir . --name MessageRepo --output ./mocks
type MessageRepo interface{}

//go:generate mockery --dir . --name ReplieRepo --output ./mocks
type ReplieRepo interface{}

//go:generate mockery --dir . --name WebUserRepo --output ./mocks
type WebUserRepo interface{}

//go:generate mockery --dir . --name SavedRepo --output ./mocks
type SavedRepo interface{}
