package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend_api/internal/model"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepo struct {
	db *DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) GetUserByID(ID int) (*model.User, error) {
	var user model.User

	err := u.db.Get(&user, "SELECT * FROM tg_user WHERE id = $1;", ID)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}
