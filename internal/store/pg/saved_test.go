package pg_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VladPetriv/scanner_backend_api/internal/model"
	"github.com/VladPetriv/scanner_backend_api/internal/store/pg"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_GetSavedMessages(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewSavedRepo(&pg.DB{DB: sqlxDB})

	data := []model.Saved{
		{ID: 1, UserID: 1, MessageID: 1},
		{ID: 2, UserID: 1, MessageID: 2},
	}

	tests := []struct {
		name           string
		mock           func()
		input          int
		want           []model.Saved
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [saved messages found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "message_id"}).
					AddRow(1, 1, 1).
					AddRow(2, 1, 2)

				mock.ExpectQuery("SELECT * FROM saved WHERE user_id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  data,
		},
		{
			name: "Error: [saved messages not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "message_id"})

				mock.ExpectQuery("SELECT * FROM saved WHERE user_id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "saved messages not found",
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM saved WHERE user_id = $1;").
					WithArgs(1).WillReturnError(fmt.Errorf("some error"))
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "failed to get saved messages by user id: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetSavedMessages(tt.input)
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

func Test_CreateSavedMessage(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewSavedRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name           string
		mock           func()
		input          *model.Saved
		want           int
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [saved message created]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1)

				mock.ExpectQuery("INSERT INTO saved(user_id, message_id) VALUES ($1, $2) RETURNING id;").
					WithArgs(1, 1).WillReturnRows(rows)
			},
			input: &model.Saved{UserID: 1, MessageID: 1},
			want:  1,
		},
		{
			name: "Error: [saved message not created]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})

				mock.ExpectQuery("INSERT INTO saved(user_id, message_id) VALUES ($1, $2) RETURNING id;").
					WithArgs(1, 1).WillReturnRows(rows)
			},
			input:          &model.Saved{UserID: 1, MessageID: 1},
			wantErr:        true,
			expectedErrMsg: "saved message not created",
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("INSERT INTO saved(user_id, message_id) VALUES ($1, $2) RETURNING id;").
					WithArgs(1, 1).WillReturnError(fmt.Errorf("some error"))
			},
			input:          &model.Saved{UserID: 1, MessageID: 1},
			wantErr:        true,
			expectedErrMsg: "failed to create saved message: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateSavedMessage(tt.input)
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

func Test_DeleteSavedMessage(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewSavedRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name           string
		mock           func()
		input          int
		want           int
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [saved message deleted]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1)

				mock.ExpectQuery("DELETE FROM saved WHERE id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  1,
		},
		{
			name: "Error: [saved message not deleted]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})

				mock.ExpectQuery("DELETE FROM saved WHERE id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "saved message not deleted",
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("DELETE FROM saved WHERE id = $1;").
					WithArgs(1).WillReturnError(fmt.Errorf("some error"))
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "failed to delete saved message by id: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.DeleteSavedMessage(tt.input)
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
