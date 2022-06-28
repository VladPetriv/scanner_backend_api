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

func (r *ReplieRepo) GetFullRepliesByMessageID(ID int) ([]model.FullReplie, error) {
	replies := make([]model.FullReplie, 0, 10)

	err := r.db.Select(
		&replies,
		`SELECT
		r.id, r.title, r.message_id as messageId 
		u.id as userId, u.fullname, u.imageurl 	
		FROM replie r 
		LEFT JOIN tg_user u ON u.id = r.user_id,
		WHERE r.message_id = $1;`,
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
