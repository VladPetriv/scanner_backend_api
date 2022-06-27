package pg

import (
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend_api/internal/model"
)

var ErrFullMessagesNotFound = errors.New("full messages not found")

type MessageRepo struct {
	db *DB
}

func NewMessageRepo(db *DB) *MessageRepo {
	return &MessageRepo{db: db}
}

func (m *MessageRepo) GetMessagesCount() (int, error) {
	var count int

	err := m.db.Get(&count, "SELECT COUNT(*) FROM message;")
	if err != nil {
		return count, fmt.Errorf("failed to get messages count: %w", err)
	}

	return count, nil
}

func (m *MessageRepo) GetMessagesCountByChannelID(ID int) (int, error) {
	var count int

	err := m.db.Get(&count, "SELECT COUNT(*) FROM message WHERE channel_id = $1;", ID)
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
	if err != nil {
		return nil, fmt.Errorf("failed to get full message by id: %w", err)
	}

	if message.Title == "" {
		return nil, ErrFullMessagesNotFound
	}

	return &message, nil
}
