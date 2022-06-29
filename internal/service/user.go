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

func (u *UserDBService) GetUserByID(ID int) (*model.User, error) {
	user, err := u.store.User.GetUserByID(ID)
	if err != nil {
		return nil, fmt.Errorf("[User] srv.GetUserByID error: %w", err)
	}

	return user, nil
}
