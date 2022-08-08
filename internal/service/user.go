package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend_api/internal/model"
	"github.com/VladPetriv/scanner_backend_api/internal/store"
)

type UserDBService struct {
	store *store.Store
}

func NewUserService(store *store.Store) *UserDBService {
	return &UserDBService{store: store}
}

func (u *UserDBService) CreateUser(user *model.UserDTO) (int, error) {
	candidate, _ := u.GetUserByUsername(user.Username)
	if candidate != nil {
		return candidate.ID, nil
	}

	id, err := u.store.User.CreateUser(user)
	if err != nil {
		return 0, fmt.Errorf("[User] srv.CreateUser error: %w", err)
	}

	return id, nil
}

func (u *UserDBService) GetUserByUsername(username string) (*model.User, error) {
	user, err := u.store.User.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("[User] srv.GetUserByUsername error: %w", err)
	}

	return user, nil
}

func (u *UserDBService) GetUserByID(ID int) (*model.User, error) {
	user, err := u.store.User.GetUserByID(ID)
	if err != nil {
		return nil, fmt.Errorf("[User] srv.GetUserByID error: %w", err)
	}

	return user, nil
}
