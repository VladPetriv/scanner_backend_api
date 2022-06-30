package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend_api/internal/model"
)

var (
	ErrSavedMessagesNotFound  = errors.New("saved messages not found")
	ErrSavedMessageNotCreated = errors.New("saved message not created")
	ErrSavedMessageNotDeleted = errors.New("saved message not deleted")
)

type SavedRepo struct {
	db *DB
}

func NewSavedRepo(db *DB) *SavedRepo {
	return &SavedRepo{db: db}
}

func (s *SavedRepo) GetSavedMessages(ID int) ([]model.Saved, error) {
	savedMessages := make([]model.Saved, 0, 10)

	err := s.db.Select(&savedMessages, "SELECT * FROM saved WHERE user_id = $1;", ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get saved messages by user id: %w", err)
	}

	if len(savedMessages) == 0 {
		return nil, ErrSavedMessagesNotFound
	}

	return savedMessages, nil
}

func (s *SavedRepo) CreateSavedMessage(savedMessage *model.Saved) (int, error) {
	var id int

	row := s.db.QueryRow(
		"INSERT INTO saved(user_id, message_id) VALUES ($1, $2) RETURNING id;",
		savedMessage.UserID, savedMessage.MessageID,
	)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrSavedMessageNotCreated
		}

		return 0, fmt.Errorf("failed to create saved message: %w", err)
	}

	return id, nil
}

func (s *SavedRepo) DeleteSavedMessage(ID int) (int, error) {
	var id int

	row := s.db.QueryRow("DELETE FROM saved WHERE id = $1;", ID)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrSavedMessageNotDeleted
		}

		return 0, fmt.Errorf("failed to delete saved message by id: %w", err)
	}

	return id, nil
}
