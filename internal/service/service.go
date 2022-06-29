package service

import "github.com/VladPetriv/scanner_backend_api/internal/model"

//go:generate mockery --dir . --name ChannelService --output ./mocks
type ChannelService interface {
	GetChannelsCount() (int, error)
	GetChannelsByPage(page int) ([]model.Channel, error)
	GetChannelByName(name string) (*model.Channel, error)
}

//go:generate mockery --dir . --name MessageService --output ./mocks
type MessageService interface {
	GetMessagesCount() (int, error)
	GetMessagesCountByChannelID(ID int) (int, error)
	GetFullMessageByID(ID int) (*model.FullMessage, error)
	GetFullMessagesByPage(offset int) ([]model.FullMessage, error)
	GetFullMessagesByChannelIDAndPage(ID, offset int) ([]model.FullMessage, error)
	GetFullMessagesByUserID(ID int) ([]model.FullMessage, error)
}

//go:generate mockery --dir . --name ReplieService --output ./mocks
type ReplieService interface {
	GetFullRepliesByMessageID(ID int) ([]model.FullReplie, error)
}

//go:generate mockery --dir . --name UserService --output ./mocks
type UserService interface {
	GetUserByID(ID int) (*model.User, error)
}

//go:generate mockery --dir . --name WebUserService --output ./mocks
type WebUserService interface {
	GetWebUserByEmail(email string) (*model.WebUser, error)
	CreateWebUser(user *model.WebUser) error
}

//go:generate mockery --dir . --name SavedService --output ./mocks
type SavedService interface{}
