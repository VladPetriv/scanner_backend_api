package pg

import (
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend_api/internal/model"
)

var ErrFullRepliesNotFound = errors.New("full replies not found")

type ReplieRepo struct {
	db *DB
}

func NewReplieRepo(db *DB) *ReplieRepo {
	return &ReplieRepo{db: db}
}

func (r *ReplieRepo) CreateReplie(replie *model.ReplieDTO) error {
	_, err := r.db.Exec(
		"INSERT INTO replie(message_id, user_id, title, imageurl) VALUES ($1, $2, $3, $4);",
		replie.MessageID, replie.UserID, replie.Title, replie.ImageURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create replie: %w", err)
	}

	return nil
}

func (r *ReplieRepo) GetFullRepliesByMessageID(ID int) ([]model.FullReplie, error) {
	replies := make([]model.FullReplie, 0, 10)

	err := r.db.Select(
		&replies,
		`SELECT
		r.id, r.title, r.message_id, r.imageurl,
		u.id as userId, u.fullname, u.imageurl AS userimageurl
		FROM replie r 
		LEFT JOIN tg_user u ON u.id = r.user_id
		WHERE r.message_id = $1
		ORDER BY r.id DESC NULLS LAST;`,
		ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get full replies by message ID: %w", err)
	}

	if len(replies) == 0 {
		return nil, ErrFullRepliesNotFound
	}

	return replies, nil
}
