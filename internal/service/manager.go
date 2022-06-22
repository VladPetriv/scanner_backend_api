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
}

func New(store *store.Store) (*Manager, error) {
	if store == nil {
		return nil, ErrNoStore
	}

	return &Manager{}, nil
}
