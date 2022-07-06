package pg_test

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VladPetriv/scanner_backend_api/internal/model"
	"github.com/VladPetriv/scanner_backend_api/internal/store/pg"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_GetMessagesCount(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name           string
		mock           func()
		want           int
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [messages count found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).
					AddRow(10)

				mock.ExpectQuery("SELECT COUNT(*) FROM message;").
					WillReturnRows(rows)
			},
			want: 10,
		},
		{
			name: "Error: [messages count not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"})

				mock.ExpectQuery("SELECT COUNT(*) FROM message;").
					WillReturnRows(rows)
			},
			wantErr:        true,
			expectedErrMsg: "messages count not found",
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("SELECT COUNT(*) FROM message;").
					WillReturnError(fmt.Errorf("some error"))
			},
			wantErr:        true,
			expectedErrMsg: "failed to get messages count: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetMessagesCount()
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_GetMessagesCountByChannelID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name           string
		mock           func()
		input          int
		want           int
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [messages count found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).
					AddRow(10)

				mock.ExpectQuery("SELECT COUNT(*) FROM message WHERE channel_id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  10,
		},
		{
			name: "Error: [messages count not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"})

				mock.ExpectQuery("SELECT COUNT(*) FROM message WHERE channel_id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "messages count not found",
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("SELECT COUNT(*) FROM message WHERE channel_id = $1;").
					WithArgs(1).WillReturnError(fmt.Errorf("some error"))
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "failed to get messages count by channel id: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetMessagesCountByChannelID(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_GetFullMessageByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: sqlxDB})

	data := &model.FullMessage{
		ID:              1,
		Title:           "test",
		MessageURL:      "test.tg",
		MessageImageURL: "test.jpg",
		ChannelName:     "testc",
		ChannelTitle:    "testc testc",
		ChannelImageURL: "testc.jpg",
		UserID:          1,
		UserFullname:    "testu testu",
		UserImageURL:    "testu.jpg",
		RepliesCount:    0,
	}

	tests := []struct {
		name           string
		mock           func()
		input          int
		want           *model.FullMessage
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [full message by ID found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"messageid", "messagetitle", "messageurl", "messageimageurl", "channelname", "channeltitle", "channelimageurl", "userid", "userfullname", "userimageurl", "count"}).
					AddRow(1, "test", "test.tg", "test.jpg", "testc", "testc testc", "testc.jpg", 1, "testu testu", "testu.jpg", 0)

				mock.ExpectQuery(
					`SELECT 
					m.id as messageId, m.title as messageTitle, m.message_url as messageUrl, m.imageurl as messageImageUrl, 
					c.name as channelName, c.title as channelTitle, c.imageurl as channelImageUrl, 
					u.id as userId, u.fullname as userFullname, u.imageurl as userImageUrl,
					(SELECT COUNT(*) FROM replie WHERE message_id = m.id)
					FROM message m
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON m.user_id = u.id 
					WHERE m.id = $1;`,
				).
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  data,
		},
		{
			name: "Error: [full message by ID not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"messageid", "messagetitle", "messageurl", "messageimageurl", "channelname", "channeltitle", "channelimageurl", "userid", "userfullname", "userimageurl", "count"})

				mock.ExpectQuery(
					`SELECT 
					m.id as messageId, m.title as messageTitle, m.message_url as messageUrl, m.imageurl as messageImageUrl, 
					c.name as channelName, c.title as channelTitle, c.imageurl as channelImageUrl, 
					u.id as userId, u.fullname as userFullname, u.imageurl as userImageUrl,
					(SELECT COUNT(*) FROM replie WHERE message_id = m.id)
					FROM message m
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON m.user_id = u.id 
					WHERE m.id = $1;`,
				).
					WithArgs(1).WillReturnRows(rows)
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "full message not found",
		},
		{
			name: "Error: [some error]",
			mock: func() {
				mock.ExpectQuery(
					`SELECT 
					m.id as messageId, m.title as messageTitle, m.message_url as messageUrl, m.imageurl as messageImageUrl, 
					c.name as channelName, c.title as channelTitle, c.imageurl as channelImageUrl, 
					u.id as userId, u.fullname as userFullname, u.imageurl as userImageUrl,
					(SELECT COUNT(*) FROM replie WHERE message_id = m.id)
					FROM message m
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON m.user_id = u.id 
					WHERE m.id = $1;`,
				).
					WithArgs(1).WillReturnError(fmt.Errorf("some error"))
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "failed to get full message by id: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetFullMessageByID(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_GetFullMessagesByPage(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: sqlxDB})

	var data []model.FullMessage

	file, err := os.Open("./fake_data/fullmessages.json")
	if err != nil {
		panic(err)
	}

	bytes, _ := io.ReadAll(file)

	json.Unmarshal(bytes, &data)

	tests := []struct {
		name           string
		mock           func()
		input          int
		want           []model.FullMessage
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [full messages with offset 0 found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"messageid", "messagetitle", "messageurl", "messageimageurl", "channelname", "channeltitle", "channelimageurl", "userid", "userfullname", "userimageurl", "count"}).
					AddRow(1, "test1", "test1.tg", "test1.jpg", "test1c", "test1c testc", "test1c.jpg", 1, "test1u testu", "test1u.jpg", 0).
					AddRow(2, "test2", "test2.tg", "test2.jpg", "test2c", "test2c testc", "test2c.jpg", 2, "test2u testu", "test2u.jpg", 0).
					AddRow(3, "test3", "test3.tg", "test3.jpg", "test3c", "test3c testc", "test3c.jpg", 3, "test3u testu", "test3u.jpg", 0).
					AddRow(4, "test4", "test4.tg", "test4.jpg", "test4c", "test4c testc", "test4c.jpg", 4, "test4u testu", "test4u.jpg", 0).
					AddRow(5, "test5", "test5.tg", "test5.jpg", "test5c", "test5c testc", "test5c.jpg", 5, "test5u testu", "test5u.jpg", 0).
					AddRow(6, "test6", "test6.tg", "test6.jpg", "test6c", "test6c testc", "test6c.jpg", 6, "test6u testu", "test6u.jpg", 0).
					AddRow(7, "test7", "test7.tg", "test7.jpg", "test7c", "test7c testc", "test7c.jpg", 7, "test7u testu", "test7u.jpg", 0).
					AddRow(8, "test8", "test8.tg", "test8.jpg", "test8c", "test8c testc", "test8c.jpg", 8, "test8u testu", "test8u.jpg", 0).
					AddRow(9, "test9", "test9.tg", "test9.jpg", "test9c", "test9c testc", "test9c.jpg", 9, "test9u testu", "test9u.jpg", 0).
					AddRow(10, "test10", "test10.tg", "test10.jpg", "test10c", "test10c testc", "test10c.jpg", 10, "test10u testu", "test10u.jpg", 0)

				mock.ExpectQuery(
					`SELECT 
					m.id as messageId, m.title as messageTitle, m.message_url as messageUrl, m.imageurl as messageImageUrl, 
					c.name as channelName, c.title as channelTitle, c.imageurl as channelImageUrl, 
					u.id as userId, u.fullname as userFullname, u.imageurl as userImageUrl,
					(SELECT COUNT(*) FROM replie WHERE message_id = m.id)
					FROM message m
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON m.user_id = u.id 
					ORDER BY m.id DESC NULLS LAST OFFSET $1 LIMIT 10;`,
				).
					WithArgs(0).WillReturnRows(rows)
			},
			input: 0,
			want:  data[:10],
		},
		{
			name: "Ok: [full messages with offset 10 found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"messageid", "messagetitle", "messageurl", "messageimageurl", "channelname", "channeltitle", "channelimageurl", "userid", "userfullname", "userimageurl", "count"}).
					AddRow(11, "test11", "test11.tg", "test11.jpg", "test11c", "test11c testc", "test11c.jpg", 11, "test11u testu", "test11u.jpg", 0).
					AddRow(12, "test12", "test12.tg", "test12.jpg", "test12c", "test12c testc", "test12c.jpg", 12, "test12u testu", "test12u.jpg", 0).
					AddRow(13, "test13", "test13.tg", "test13.jpg", "test13c", "test13c testc", "test13c.jpg", 13, "test13u testu", "test13u.jpg", 0).
					AddRow(14, "test14", "test14.tg", "test14.jpg", "test14c", "test14c testc", "test14c.jpg", 14, "test14u testu", "test14u.jpg", 0).
					AddRow(15, "test15", "test15.tg", "test15.jpg", "test15c", "test15c testc", "test15c.jpg", 15, "test15u testu", "test15u.jpg", 0).
					AddRow(16, "test16", "test16.tg", "test16.jpg", "test16c", "test16c testc", "test16c.jpg", 16, "test16u testu", "test16u.jpg", 0).
					AddRow(17, "test17", "test17.tg", "test17.jpg", "test17c", "test17c testc", "test17c.jpg", 17, "test17u testu", "test17u.jpg", 0).
					AddRow(18, "test18", "test18.tg", "test18.jpg", "test18c", "test18c testc", "test18c.jpg", 18, "test18u testu", "test18u.jpg", 0).
					AddRow(19, "test19", "test19.tg", "test19.jpg", "test19c", "test19c testc", "test19c.jpg", 19, "test19u testu", "test19u.jpg", 0).
					AddRow(20, "test20", "test20.tg", "test20.jpg", "test20c", "test20c testc", "test20c.jpg", 20, "test20u testu", "test20u.jpg", 0)

				mock.ExpectQuery(
					`SELECT 
					m.id as messageId, m.title as messageTitle, m.message_url as messageUrl, m.imageurl as messageImageUrl, 
					c.name as channelName, c.title as channelTitle, c.imageurl as channelImageUrl, 
					u.id as userId, u.fullname as userFullname, u.imageurl as userImageUrl,
					(SELECT COUNT(*) FROM replie WHERE message_id = m.id)
					FROM message m
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON m.user_id = u.id 
					ORDER BY m.id DESC NULLS LAST OFFSET $1 LIMIT 10;`,
				).WithArgs(10).WillReturnRows(rows)
			},
			input: 10,
			want:  data[10:],
		},
		{
			name: "Error: [full messages not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"messageid", "messagetitle", "messageurl", "messageimageurl", "channelname", "channeltitle", "channelimageurl", "userid", "userfullname", "userimageurl", "count"})

				mock.ExpectQuery(
					`SELECT 
					m.id as messageId, m.title as messageTitle, m.message_url as messageUrl, m.imageurl as messageImageUrl, 
					c.name as channelName, c.title as channelTitle, c.imageurl as channelImageUrl, 
					u.id as userId, u.fullname as userFullname, u.imageurl as userImageUrl,
					(SELECT COUNT(*) FROM replie WHERE message_id = m.id)
					FROM message m
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON m.user_id = u.id 
					ORDER BY m.id DESC NULLS LAST OFFSET $1 LIMIT 10;`,
				).WithArgs(0).WillReturnRows(rows)
			},
			input:          0,
			wantErr:        true,
			expectedErrMsg: "full messages not found",
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery(
					`SELECT 
					m.id as messageId, m.title as messageTitle, m.message_url as messageUrl, m.imageurl as messageImageUrl, 
					c.name as channelName, c.title as channelTitle, c.imageurl as channelImageUrl, 
					u.id as userId, u.fullname as userFullname, u.imageurl as userImageUrl,
					(SELECT COUNT(*) FROM replie WHERE message_id = m.id)
					FROM message m
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON m.user_id = u.id 
					ORDER BY m.id DESC NULLS LAST OFFSET $1 LIMIT 10;`,
				).WithArgs(0).WillReturnError(fmt.Errorf("some error"))
			},
			input:          0,
			wantErr:        true,
			expectedErrMsg: "failed to get full messages by page: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetFullMessagesByPage(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_GetFullMessagesByChannelIDAndPage(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: sqlxDB})

	var data []model.FullMessage

	file, err := os.Open("./fake_data/fullmessagesbychannel.json")
	if err != nil {
		panic(err)
	}

	bytes, _ := io.ReadAll(file)

	json.Unmarshal(bytes, &data)

	tests := []struct {
		name           string
		mock           func()
		ID             int
		offset         int
		want           []model.FullMessage
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [full messages by channel ID with offset 0 found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"messageid", "messagetitle", "messageurl", "messageimageurl", "channelname", "channeltitle", "channelimageurl", "userid", "userfullname", "userimageurl", "count"}).
					AddRow(1, "test1", "test1.tg", "test1.jpg", "testc", "testc testc", "testc.jpg", 1, "test1u testu", "test1u.jpg", 0).
					AddRow(2, "test2", "test2.tg", "test2.jpg", "testc", "testc testc", "testc.jpg", 2, "test2u testu", "test2u.jpg", 0).
					AddRow(3, "test3", "test3.tg", "test3.jpg", "testc", "testc testc", "testc.jpg", 3, "test3u testu", "test3u.jpg", 0).
					AddRow(4, "test4", "test4.tg", "test4.jpg", "testc", "testc testc", "testc.jpg", 4, "test4u testu", "test4u.jpg", 0).
					AddRow(5, "test5", "test5.tg", "test5.jpg", "testc", "testc testc", "testc.jpg", 5, "test5u testu", "test5u.jpg", 0).
					AddRow(6, "test6", "test6.tg", "test6.jpg", "testc", "testc testc", "testc.jpg", 6, "test6u testu", "test6u.jpg", 0).
					AddRow(7, "test7", "test7.tg", "test7.jpg", "testc", "testc testc", "testc.jpg", 7, "test7u testu", "test7u.jpg", 0).
					AddRow(8, "test8", "test8.tg", "test8.jpg", "testc", "testc testc", "testc.jpg", 8, "test8u testu", "test8u.jpg", 0).
					AddRow(9, "test9", "test9.tg", "test9.jpg", "testc", "testc testc", "testc.jpg", 9, "test9u testu", "test9u.jpg", 0).
					AddRow(10, "test10", "test10.tg", "test10.jpg", "testc", "testc testc", "testc.jpg", 10, "test10u testu", "test10u.jpg", 0)

				mock.ExpectQuery(
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
				).
					WithArgs(1, 0).WillReturnRows(rows)
			},
			ID:     1,
			offset: 0,
			want:   data[:10],
		},
		{
			name: "Ok: [full messages by channel ID with offset 10 found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"messageid", "messagetitle", "messageurl", "messageimageurl", "channelname", "channeltitle", "channelimageurl", "userid", "userfullname", "userimageurl", "count"}).
					AddRow(11, "test11", "test11.tg", "test11.jpg", "testc", "testc testc", "testc.jpg", 11, "test11u testu", "test11u.jpg", 0).
					AddRow(12, "test12", "test12.tg", "test12.jpg", "testc", "testc testc", "testc.jpg", 12, "test12u testu", "test12u.jpg", 0).
					AddRow(13, "test13", "test13.tg", "test13.jpg", "testc", "testc testc", "testc.jpg", 13, "test13u testu", "test13u.jpg", 0).
					AddRow(14, "test14", "test14.tg", "test14.jpg", "testc", "testc testc", "testc.jpg", 14, "test14u testu", "test14u.jpg", 0).
					AddRow(15, "test15", "test15.tg", "test15.jpg", "testc", "testc testc", "testc.jpg", 15, "test15u testu", "test15u.jpg", 0).
					AddRow(16, "test16", "test16.tg", "test16.jpg", "testc", "testc testc", "testc.jpg", 16, "test16u testu", "test16u.jpg", 0).
					AddRow(17, "test17", "test17.tg", "test17.jpg", "testc", "testc testc", "testc.jpg", 17, "test17u testu", "test17u.jpg", 0).
					AddRow(18, "test18", "test18.tg", "test18.jpg", "testc", "testc testc", "testc.jpg", 18, "test18u testu", "test18u.jpg", 0).
					AddRow(19, "test19", "test19.tg", "test19.jpg", "testc", "testc testc", "testc.jpg", 19, "test19u testu", "test19u.jpg", 0).
					AddRow(20, "test20", "test20.tg", "test20.jpg", "testc", "testc testc", "testc.jpg", 20, "test20u testu", "test20u.jpg", 0)

				mock.ExpectQuery(
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
				).WithArgs(1, 10).WillReturnRows(rows)
			},
			ID:     1,
			offset: 10,
			want:   data[10:],
		},
		{
			name: "Error: [full messages by channel ID not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"messageid", "messagetitle", "messageurl", "messageimageurl", "channelname", "channeltitle", "channelimageurl", "userid", "userfullname", "userimageurl", "count"})

				mock.ExpectQuery(
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
				).WithArgs(1, 0).WillReturnRows(rows)
			},
			ID:             1,
			offset:         0,
			wantErr:        true,
			expectedErrMsg: "full messages not found",
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery(
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
				).WithArgs(1, 0).WillReturnError(fmt.Errorf("some error"))
			},
			ID:             1,
			offset:         0,
			wantErr:        true,
			expectedErrMsg: "failed to get full messages by channel ID and page: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetFullMessagesByChannelIDAndPage(tt.ID, tt.offset)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_GetFullMessagesByUserID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: sqlxDB})

	var data []model.FullMessage

	file, err := os.Open("./fake_data/fullmessagesbyuser.json")
	if err != nil {
		panic(err)
	}

	bytes, _ := io.ReadAll(file)

	json.Unmarshal(bytes, &data)

	tests := []struct {
		name           string
		mock           func()
		ID             int
		offset         int
		want           []model.FullMessage
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [full messages by user ID found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"messageid", "messagetitle", "messageurl", "messageimageurl", "channelname", "channeltitle", "channelimageurl", "userid", "userfullname", "userimageurl", "count"}).
					AddRow(1, "test1", "test1.tg", "test1.jpg", "test1c", "test1c testc", "test1c.jpg", 1, "testu testu", "testu.jpg", 0).
					AddRow(2, "test2", "test2.tg", "test2.jpg", "test2c", "test2c testc", "test2c.jpg", 1, "testu testu", "testu.jpg", 0).
					AddRow(3, "test3", "test3.tg", "test3.jpg", "test3c", "test3c testc", "test3c.jpg", 1, "testu testu", "testu.jpg", 0).
					AddRow(4, "test4", "test4.tg", "test4.jpg", "test4c", "test4c testc", "test4c.jpg", 1, "testu testu", "testu.jpg", 0).
					AddRow(5, "test5", "test5.tg", "test5.jpg", "test5c", "test5c testc", "test5c.jpg", 1, "testu testu", "testu.jpg", 0).
					AddRow(6, "test6", "test6.tg", "test6.jpg", "test6c", "test6c testc", "test6c.jpg", 1, "testu testu", "testu.jpg", 0).
					AddRow(7, "test7", "test7.tg", "test7.jpg", "test7c", "test7c testc", "test7c.jpg", 1, "testu testu", "testu.jpg", 0).
					AddRow(8, "test8", "test8.tg", "test8.jpg", "test8c", "test8c testc", "test8c.jpg", 1, "testu testu", "testu.jpg", 0).
					AddRow(9, "test9", "test9.tg", "test9.jpg", "test9c", "test9c testc", "test9c.jpg", 1, "testu testu", "testu.jpg", 0).
					AddRow(10, "test10", "test10.tg", "test10.jpg", "test10c", "test10c testc", "test10c.jpg", 1, "testu testu", "testu.jpg", 0)

				mock.ExpectQuery(
					`SELECT 
					m.id as messageId, m.title as messageTitle, m.message_url as messageUrl, m.imageurl as messageImageUrl, 
					c.name as channelName, c.title as channelTitle, c.imageurl as channelImageUrl, 
					u.id as userId, u.fullname as userFullname, u.imageurl as userImageUrl,
					(SELECT COUNT(*) FROM replie WHERE message_id = m.id)
					FROM message m
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON m.user_id = u.id 
					WHERE m.user_id = $1;`,
				).
					WithArgs(1).WillReturnRows(rows)
			},
			ID:   1,
			want: data[:10],
		},
		{
			name: "Ok: [full messages by user ID found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"messageid", "messagetitle", "messageurl", "messageimageurl", "channelname", "channeltitle", "channelimageurl", "userid", "userfullname", "userimageurl", "count"}).
					AddRow(11, "test11", "test11.tg", "test11.jpg", "test11c", "test11c testc", "test11c.jpg", 2, "test2u testu", "test2u.jpg", 0).
					AddRow(12, "test12", "test12.tg", "test12.jpg", "test12c", "test12c testc", "test12c.jpg", 2, "test2u testu", "test2u.jpg", 0).
					AddRow(13, "test13", "test13.tg", "test13.jpg", "test13c", "test13c testc", "test13c.jpg", 2, "test2u testu", "test2u.jpg", 0).
					AddRow(14, "test14", "test14.tg", "test14.jpg", "test14c", "test14c testc", "test14c.jpg", 2, "test2u testu", "test2u.jpg", 0).
					AddRow(15, "test15", "test15.tg", "test15.jpg", "test15c", "test15c testc", "test15c.jpg", 2, "test2u testu", "test2u.jpg", 0).
					AddRow(16, "test16", "test16.tg", "test16.jpg", "test16c", "test16c testc", "test16c.jpg", 2, "test2u testu", "test2u.jpg", 0).
					AddRow(17, "test17", "test17.tg", "test17.jpg", "test17c", "test17c testc", "test17c.jpg", 2, "test2u testu", "test2u.jpg", 0).
					AddRow(18, "test18", "test18.tg", "test18.jpg", "test18c", "test18c testc", "test18c.jpg", 2, "test2u testu", "test2u.jpg", 0).
					AddRow(19, "test19", "test19.tg", "test19.jpg", "test19c", "test19c testc", "test19c.jpg", 2, "test2u testu", "test2u.jpg", 0).
					AddRow(20, "test20", "test20.tg", "test20.jpg", "test20c", "test20c testc", "test20c.jpg", 2, "test2u testu", "test2u.jpg", 0)

				mock.ExpectQuery(
					`SELECT 
					m.id as messageId, m.title as messageTitle, m.message_url as messageUrl, m.imageurl as messageImageUrl, 
					c.name as channelName, c.title as channelTitle, c.imageurl as channelImageUrl, 
					u.id as userId, u.fullname as userFullname, u.imageurl as userImageUrl,
					(SELECT COUNT(*) FROM replie WHERE message_id = m.id)
					FROM message m
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON m.user_id = u.id 
					WHERE m.user_id = $1;`,
				).WithArgs(2).WillReturnRows(rows)
			},
			ID:   2,
			want: data[10:],
		},
		{
			name: "Error: [full messages by user ID not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"messageid", "messagetitle", "messageurl", "messageimageurl", "channelname", "channeltitle", "channelimageurl", "userid", "userfullname", "userimageurl", "count"})

				mock.ExpectQuery(
					`SELECT 
					m.id as messageId, m.title as messageTitle, m.message_url as messageUrl, m.imageurl as messageImageUrl, 
					c.name as channelName, c.title as channelTitle, c.imageurl as channelImageUrl, 
					u.id as userId, u.fullname as userFullname, u.imageurl as userImageUrl,
					(SELECT COUNT(*) FROM replie WHERE message_id = m.id)
					FROM message m
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON m.user_id = u.id 
					WHERE m.user_id = $1;`,
				).WithArgs(1).WillReturnRows(rows)
			},
			ID:             1,
			wantErr:        true,
			expectedErrMsg: "full messages not found",
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery(
					`SELECT 
					m.id as messageId, m.title as messageTitle, m.message_url as messageUrl, m.imageurl as messageImageUrl, 
					c.name as channelName, c.title as channelTitle, c.imageurl as channelImageUrl, 
					u.id as userId, u.fullname as userFullname, u.imageurl as userImageUrl,
					(SELECT COUNT(*) FROM replie WHERE message_id = m.id)
					FROM message m
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON m.user_id = u.id 
					WHERE m.user_id = $1;`,
				).WithArgs(1).WillReturnError(fmt.Errorf("some error"))
			},
			ID:             1,
			wantErr:        true,
			expectedErrMsg: "failed to get full messages by user ID: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetFullMessagesByUserID(tt.ID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
