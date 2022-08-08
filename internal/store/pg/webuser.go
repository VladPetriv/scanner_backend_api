package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend_api/internal/model"
)

var (
	ErrWebUserNotFound   = errors.New("web user not found")
	ErrWebUserNotCreated = errors.New("web user not created")
)

type WebUserRepo struct {
	db *DB
}

func NewWebUserRepo(db *DB) *WebUserRepo {
	return &WebUserRepo{db: db}
}

func (w *WebUserRepo) GetWebUserByEmail(email string) (*model.WebUser, error) {
	var user model.WebUser

	err := w.db.Get(&user, "SELECT * FROM web_user WHERE email = $1;", email)
	if err == sql.ErrNoRows {
		return nil, ErrWebUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get web user by email: %w", err)
	}

	return &user, nil
}

func (w *WebUserRepo) CreateWebUser(user *model.WebUser) (int, error) {
	var id int

	row := w.db.QueryRow("INSERT INTO web_user(email, password) VALUES ($1, $2) RETURNING id;", user.Email, user.Password)

	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrWebUserNotCreated
		}

		return 0, fmt.Errorf("failed to create web user:%w", err)
	}

	return id, nil
}
