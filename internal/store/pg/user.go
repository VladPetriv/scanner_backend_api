package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend_api/internal/model"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUserNotCreated = errors.New("user not created")
)

type UserRepo struct {
	db *DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) CreateUser(user *model.UserDTO) (int, error) {
	var id int

	row := u.db.QueryRow(
		"INSERT INTO tg_user(username, fullname, imageurl) VALUES ($1, $2, $3) RETURNING id;",
		user.Username, user.Fullname, user.ImageURL,
	)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrUserNotCreated
		}

		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (u *UserRepo) GetUserByUsername(username string) (*model.User, error) {
	var user model.User

	err := u.db.Get(&user, "SELECT * FROM tg_user WHERE username = $1;", username)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
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
