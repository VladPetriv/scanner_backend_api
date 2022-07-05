package service

import (
	"errors"

	"github.com/VladPetriv/scanner_backend_api/internal/store"
)

var ErrNoStore = errors.New("no store provided")

type Manager struct {
	Channel ChannelService
	Message MessageService
	Replie  ReplieService
	User    UserService
	WebUser WebUserService
	Saved   SavedService
	Jwt     JwtService
}

func New(store *store.Store, secretJWTKey string) (*Manager, error) {
	if store == nil {
		return nil, ErrNoStore
	}

	return &Manager{
		Channel: NewChannelService(store),
		Message: NewMessageService(store),
		Replie:  NewReplieService(store),
		User:    NewUserService(store),
		WebUser: NewWebUserService(store),
		Saved:   NewSavedService(store),
		Jwt:     NewJwtService(secretJWTKey),
	}, nil
}
