package service

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend_api/internal/model"
	"github.com/VladPetriv/scanner_backend_api/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type WebUserDBService struct {
	store *store.Store
}

func NewWebUserService(store *store.Store) *WebUserDBService {
	return &WebUserDBService{store: store}
}

func (w *WebUserDBService) GetWebUserByEmail(email string) (*model.WebUser, error) {
	user, err := w.store.WebUser.GetWebUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("[WebUser] srv.GetWebUserByEmail error: %w", err)
	}

	return user, nil
}

func (w *WebUserDBService) CreateWebUser(user *model.WebUser) error {
	user.Password, _ = w.HashPassword(user.Password)

	_, err := w.store.WebUser.CreateWebUser(user)
	if err != nil {
		return fmt.Errorf("[WebUser] srv.CreateWebUser error: %w", err)
	}

	return nil
}

func (w *WebUserDBService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(bytes), nil
}

func (w *WebUserDBService) ComparePassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}
