package service

import "github.com/VladPetriv/scanner_backend_api/internal/model"

//go:generate mockery --dir . --name ChannelService --output ./mocks
type ChannelService interface {
	CreateChannel(channel *model.ChannelDTO) error
	GetChannelsCount() (int, error)
	GetChannelsByPage(page int) ([]model.Channel, error)
	GetChannelByName(name string) (*model.Channel, error)
}

//go:generate mockery --dir . --name MessageService --output ./mocks
type MessageService interface {
	CreateMessage(message *model.MessageDTO) (int, error)
	GetMessagesCount() (int, error)
	GetMessagesCountByChannelID(ID int) (int, error)
	GetFullMessageByID(ID int) (*model.FullMessage, error)
	GetFullMessagesByPage(offset int) ([]model.FullMessage, error)
	GetFullMessagesByChannelIDAndPage(ID, offset int) ([]model.FullMessage, error)
	GetFullMessagesByUserID(ID int) ([]model.FullMessage, error)
}

//go:generate mockery --dir . --name ReplieService --output ./mocks
type ReplieService interface {
	CreateReplie(replie *model.ReplieDTO) error
	GetFullRepliesByMessageID(ID int) ([]model.FullReplie, error)
}

//go:generate mockery --dir . --name UserService --output ./mocks
type UserService interface {
	CreateUser(user *model.UserDTO) (int, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByID(ID int) (*model.User, error)
}

//go:generate mockery --dir . --name WebUserService --output ./mocks
type WebUserService interface {
	GetWebUserByEmail(email string) (*model.WebUser, error)
	CreateWebUser(user *model.WebUser) error
	HashPassword(password string) (string, error)
	ComparePassword(password, HashPassword string) bool
}

//go:generate mockery --dir . --name SavedService --output ./mocks
type SavedService interface {
	GetSavedMessages(ID int) ([]model.Saved, error)
	CreateSavedMessage(savedMessage *model.Saved) error
	DeleteSavedMessage(ID int) error
}

//go:generate mockery --dir . --name JwtService --output ./mocks
type JwtService interface {
	GenerateToken(userEmail string) (string, error)
	ParseToken(accessToken string) (string, error)
}
