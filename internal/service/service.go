package service

import "github.com/VladPetriv/scanner_backend_api/internal/model"

//go:generate mockery --dir . --name ChannelService --output ./mocks
type ChannelService interface {
	GetChannelsCount() (int, error)
	GetChannelsByPage(page int) ([]model.Channel, error)
	GetChannelByName(name string) (*model.Channel, error)
}

//go:generate mockery --dir . --name MessageService --output ./mocks
type MessageService interface{}

//go:generate mockery --dir . --name ReplieService --output ./mocks
type ReplieService interface{}

//go:generate mockery --dir . --name UserService --output ./mocks
type UserService interface{}

//go:generate mockery --dir . --name WebUserService --output ./mocks
type WebUserService interface{}

//go:generate mockery --dir . --name SavedService --output ./mocks
type SavedService interface{}
