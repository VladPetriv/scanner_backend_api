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

func Test_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewUserRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name           string
		mock           func()
		input          *model.UserDTO
		want           int
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [user created]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1)

				mock.ExpectQuery("INSERT INTO tg_user(username, fullname, imageurl) VALUES ($1, $2, $3) RETURNING id;").
					WithArgs("test", "test test", "test.jpg").WillReturnRows(rows)
			},
			input: &model.UserDTO{Username: "test", Fullname: "test test", ImageURL: "test.jpg"},
			want:  1,
		},
		{
			name: "Error: [user not created]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})

				mock.ExpectQuery("INSERT INTO tg_user(username, fullname, imageurl) VALUES ($1, $2, $3) RETURNING id;").
					WithArgs("test", "test test", "test.jpg").WillReturnRows(rows)
			},
			input:          &model.UserDTO{Username: "test", Fullname: "test test", ImageURL: "test.jpg"},
			wantErr:        true,
			expectedErrMsg: "user not created",
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("INSERT INTO tg_user(username, fullname, imageurl) VALUES ($1, $2, $3) RETURNING id;").
					WithArgs("test", "test test", "test.jpg").WillReturnError(fmt.Errorf("some error"))
			},
			input:          &model.UserDTO{Username: "test", Fullname: "test test", ImageURL: "test.jpg"},
			wantErr:        true,
			expectedErrMsg: "failed to create user: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateUser(tt.input)
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

func Test_GetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewUserRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name           string
		mock           func()
		input          int
		want           *model.User
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [user by id found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "fullname", "imageurl"}).
					AddRow(1, "test", "test test", "test.jpg")

				mock.ExpectQuery("SELECT * FROM tg_user WHERE id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  &model.User{ID: 1, Username: "test", Fullname: "test test", ImageURL: "test.jpg"},
		},
		{
			name: "Error: [user by id not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "fullname", "imageurl"})

				mock.ExpectQuery("SELECT * FROM tg_user WHERE id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "user not found",
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM tg_user WHERE id = $1;").
					WithArgs(1).WillReturnError(fmt.Errorf("some error"))
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "failed to get user by id: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetUserByID(tt.input)
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
