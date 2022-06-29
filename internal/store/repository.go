package store

import "github.com/VladPetriv/scanner_backend_api/internal/model"

//go:generate mockery --dir . --name ChannelRepo --output ./mocks
type ChannelRepo interface {
	GetChannelsCount() (int, error)
	GetChannelsByPage(offset int) ([]model.Channel, error)
	GetChannelByName(name string) (*model.Channel, error)
}

//go:generate mockery --dir . --name MessageRepo --output ./mocks
type MessageRepo interface {
	GetMessagesCount() (int, error)
	GetMessagesCountByChannelID(ID int) (int, error)
	GetFullMessageByID(ID int) (*model.FullMessage, error)
	GetFullMessagesByPage(offset int) ([]model.FullMessage, error)
	GetFullMessagesByChannelIDAndPage(ID, offset int) ([]model.FullMessage, error)
	GetFullMessagesByUserID(ID int) ([]model.FullMessage, error)
}

//go:generate mockery --dir . --name ReplieRepo --output ./mocks
type ReplieRepo interface {
	GetFullRepliesByMessageID(ID int) ([]model.FullReplie, error)
}

//go:generate mockery --dir . --name UserRepo --output ./mocks
type UserRepo interface {
	GetUserByID(ID int) (*model.User, error)
}

//go:generate mockery --dir . --name WebUserRepo --output ./mocks
type WebUserRepo interface {
	GetWebUserByEmail(email string) (*model.WebUser, error)
	CreateWebUser(user *model.WebUser) (int, error)
}

//go:generate mockery --dir . --name SavedRepo --output ./mocks
type SavedRepo interface{}
