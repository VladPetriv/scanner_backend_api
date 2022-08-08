package pg_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/VladPetriv/scanner_backend_api/internal/model"
	"github.com/VladPetriv/scanner_backend_api/internal/store/pg"
)

func Test_CreateChannel(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewChannelRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name           string
		mock           func()
		input          *model.ChannelDTO
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [channel created]",
			mock: func() {
				mock.ExpectExec("INSERT INTO channel(name, title, imageurl) VALUES ($1, $2, $3);").
					WithArgs("test", "test T", "test.jpg").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: &model.ChannelDTO{Name: "test", Title: "test T", ImageURL: "test.jpg"},
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectExec("INSERT INTO channel(name, title, imageurl) VALUES ($1, $2, $3);").
					WithArgs("test", "test T", "test.jpg").WillReturnError(fmt.Errorf("some error"))
			},
			input:          &model.ChannelDTO{Name: "test", Title: "test T", ImageURL: "test.jpg"},
			wantErr:        true,
			expectedErrMsg: "failed to create channel: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.CreateChannel(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_GetChannelsCount(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewChannelRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name           string
		mock           func()
		want           int
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [channels count found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).
					AddRow(10)

				mock.ExpectQuery("SELECT COUNT(*) FROM channel;").
					WillReturnRows(rows)
			},
			want: 10,
		},
		{
			name: "Error: [channels count not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"})

				mock.ExpectQuery("SELECT COUNT(*) FROM channel;").
					WillReturnRows(rows)
			},
			wantErr:        true,
			expectedErrMsg: "channels count not found",
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("SELECT COUNT(*) FROM channel;").
					WillReturnError(fmt.Errorf("some error"))
			},
			wantErr:        true,
			expectedErrMsg: "failed to get count of channels: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetChannelsCount()
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

func Test_GetChannelsByPage(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewChannelRepo(&pg.DB{DB: sqlxDB})

	channels := []model.Channel{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5},
		{ID: 6}, {ID: 7}, {ID: 8}, {ID: 9}, {ID: 10},
		{ID: 11}, {ID: 12}, {ID: 13}, {ID: 14}, {ID: 15},
		{ID: 16}, {ID: 17}, {ID: 18}, {ID: 19}, {ID: 20},
	}

	tests := []struct {
		name           string
		mock           func()
		input          int
		want           []model.Channel
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [channels found with offset 0]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1).AddRow(2).AddRow(3).AddRow(4).AddRow(5).
					AddRow(6).AddRow(7).AddRow(8).AddRow(9).AddRow(10)

				mock.ExpectQuery("SELECT * FROM channel OFFSET $1 LIMIT 10;").
					WithArgs(0).WillReturnRows(rows)
			},
			input: 0,
			want:  channels[:10],
		},
		{
			name: "Ok: [channels found with offset 10]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(11).AddRow(12).AddRow(13).AddRow(14).AddRow(15).
					AddRow(16).AddRow(17).AddRow(18).AddRow(19).AddRow(20)

				mock.ExpectQuery("SELECT * FROM channel OFFSET $1 LIMIT 10;").
					WithArgs(10).WillReturnRows(rows)
			},
			input: 10,
			want:  channels[10:],
		},
		{
			name: "Error: [channels not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})

				mock.ExpectQuery("SELECT * FROM channel OFFSET $1 LIMIT 10;").
					WithArgs(0).WillReturnRows(rows)
			},
			input:          0,
			wantErr:        true,
			expectedErrMsg: "channels not found",
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM channel OFFSET $1 LIMIT 10;").
					WillReturnError(fmt.Errorf("some error"))
			},
			wantErr:        true,
			expectedErrMsg: "failed to get channels by page: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetChannelsByPage(tt.input)
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

func Test_GetChannelByName(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewChannelRepo(&pg.DB{DB: sqlxDB})

	channel := &model.Channel{
		ID:       1,
		Name:     "test",
		Title:    "test_test",
		ImageURL: "test.jpg",
	}

	tests := []struct {
		name           string
		mock           func()
		input          string
		want           *model.Channel
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [channel found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "imageurl"}).
					AddRow(1, "test", "test_test", "test.jpg")

				mock.ExpectQuery("SELECT * FROM channel WHERE name = $1;").
					WithArgs("test").WillReturnRows(rows)
			},
			input: "test",
			want:  channel,
		},
		{
			name: "Error: [channel not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "imageurl"})

				mock.ExpectQuery("SELECT * FROM channel WHERE name = $1;").
					WithArgs("404").WillReturnRows(rows)
			},
			input:          "404",
			wantErr:        true,
			expectedErrMsg: "channel not found",
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM channel WHERE name = $1;").
					WillReturnError(fmt.Errorf("some error"))
			},
			wantErr:        true,
			expectedErrMsg: "failed to get channel by name: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetChannelByName(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
