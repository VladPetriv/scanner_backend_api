package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend_api/internal/model"
)

var (
	ErrMessagesCountNotFound = errors.New("messages count not found")
	ErrFullMessagesNotFound  = errors.New("full messages not found")
	ErrFullMessageNotFound   = errors.New("full message not found")
	ErrMessageNotCreated     = errors.New("message not created")
)

type MessageRepo struct {
	db *DB
}

func NewMessageRepo(db *DB) *MessageRepo {
	return &MessageRepo{db: db}
}

func (m *MessageRepo) CreateMessage(message *model.MessageDTO) (int, error) {
	var id int

	row := m.db.QueryRow(
		"INSERT INTO message(channel_id, user_id, title, message_url, imageurl) VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		message.ChannelID, message.UserID, message.Title, message.MessageURL, message.ImageURL,
	)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrMessageNotCreated
		}

		return 0, fmt.Errorf("failed to create message: %w", err)
	}

	return id, nil
}

func (m *MessageRepo) GetMessagesCount() (int, error) {
	var count int

	err := m.db.Get(&count, "SELECT COUNT(*) FROM message;")
	if err == sql.ErrNoRows {
		return 0, ErrMessagesCountNotFound
	}

	if err != nil {
		return 0, fmt.Errorf("failed to get messages count: %w", err)
	}

	return count, nil
}

func (m *MessageRepo) GetMessagesCountByChannelID(ID int) (int, error) {
	var count int

	err := m.db.Get(&count, "SELECT COUNT(*) FROM message WHERE channel_id = $1;", ID)
	if err == sql.ErrNoRows {
		return 0, ErrMessagesCountNotFound
	}

	if err != nil {
		return 0, fmt.Errorf("failed to get messages count by channel id: %w", err)
	}

	return count, nil
}

func (m *MessageRepo) GetFullMessagesByPage(offset int) ([]model.FullMessage, error) {
	messages := make([]model.FullMessage, 0, 10)

	err := m.db.Select(
		&messages,
		`SELECT 
		m.id as messageId, m.title as messageTitle, m.message_url as messageUrl, m.imageurl as messageImageUrl, 
		c.name as channelName, c.title as channelTitle, c.imageurl as channelImageUrl, 
		u.id as userId, u.fullname as userFullname, u.imageurl as userImageUrl,
		(SELECT COUNT(*) FROM replie WHERE message_id = m.id)
		FROM message m
		LEFT JOIN channel c ON c.id = m.channel_id 
		LEFT JOIN tg_user u ON m.user_id = u.id 
		ORDER BY m.id DESC NULLS LAST OFFSET $1 LIMIT 10;`,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get full messages by page: %w", err)
	}

	if len(messages) == 0 {
		return nil, ErrFullMessagesNotFound
	}

	return messages, nil
}

func (m *MessageRepo) GetFullMessagesByChannelIDAndPage(ID int, offset int) ([]model.FullMessage, error) {
	messages := make([]model.FullMessage, 0, 10)

	err := m.db.Select(
		&messages,
		`SELECT 
		m.id as messageId, m.title as messageTitle, m.message_url as messageUrl, m.imageurl as messageImageUrl, 
		c.name as channelName, c.title as channelTitle, c.imageurl as channelImageUrl, 
		u.id as userId, u.fullname as userFullname, u.imageurl as userImageUrl,
		(SELECT COUNT(*) FROM replie WHERE message_id = m.id)
		FROM message m
		LEFT JOIN channel c ON c.id = m.channel_id 
		LEFT JOIN tg_user u ON m.user_id = u.id 
		WHERE m.channel_id = $1
		ORDER BY m.id DESC NULLS LAST OFFSET $2 LIMIT 10;`,
		ID,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get full messages by channel ID and page: %w", err)
	}

	if len(messages) == 0 {
		return nil, ErrFullMessagesNotFound
	}

	return messages, nil
}

func (m *MessageRepo) GetFullMessagesByUserID(ID int) ([]model.FullMessage, error) {
	messages := make([]model.FullMessage, 0, 10)

	err := m.db.Select(
		&messages,
		`SELECT 
		m.id as messageId, m.title as messageTitle, m.message_url as messageUrl, m.imageurl as messageImageUrl, 
		c.name as channelName, c.title as channelTitle, c.imageurl as channelImageUrl, 
		u.id as userId, u.fullname as userFullname, u.imageurl as userImageUrl,
		(SELECT COUNT(*) FROM replie WHERE message_id = m.id)
		FROM message m
		LEFT JOIN channel c ON c.id = m.channel_id 
		LEFT JOIN tg_user u ON m.user_id = u.id 
		WHERE m.user_id = $1;`,
		ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get full messages by user ID: %w", err)
	}

	if len(messages) == 0 {
		return nil, ErrFullMessagesNotFound
	}

	return messages, nil
}

func (m *MessageRepo) GetFullMessageByID(ID int) (*model.FullMessage, error) {
	var message model.FullMessage

	err := m.db.Get(
		&message,
		`SELECT 
		m.id as messageId, m.title as messageTitle, m.message_url as messageUrl, m.imageurl as messageImageUrl, 
		c.name as channelName, c.title as channelTitle, c.imageurl as channelImageUrl, 
		u.id as userId, u.fullname as userFullname, u.imageurl as userImageUrl,
		(SELECT COUNT(*) FROM replie WHERE message_id = m.id)
		FROM message m
		LEFT JOIN channel c ON c.id = m.channel_id 
		LEFT JOIN tg_user u ON m.user_id = u.id 
		WHERE m.id = $1;`,
		ID,
	)
	if err == sql.ErrNoRows {
		return nil, ErrFullMessageNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get full message by id: %w", err)
	}

	return &message, nil
}
