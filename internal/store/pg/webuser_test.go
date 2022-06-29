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

func Test_CreateWebUser(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewWebUserRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name    string
		mock    func()
		input   *model.WebUser
		want    int
		wantErr bool
	}{
		{
			name: "Ok: [web user created]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1)

				mock.ExpectQuery("INSERT INTO web_user(email, password) VALUES ($1, $2) RETURNING id;").
					WithArgs("test@test.com", "test_pswd").WillReturnRows(rows)
			},
			input: &model.WebUser{Email: "test@test.com", Password: "test_pswd"},
			want:  1,
		},
		{
			name: "Error: [no rows in result set]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})

				mock.ExpectQuery("INSERT INTO web_user(email, password) VALUES ($1, $2) RETURNING id;").
					WithArgs("test@test.com", "test_pswd").WillReturnRows(rows)
			},
			input: &model.WebUser{Email: "test@test.com", Password: "test_pswd"},

			wantErr: true,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("INSERT INTO web_user(email, password) VALUES ($1, $2) RETURNING id;").
					WithArgs("test@test.com", "test_pswd").WillReturnError(fmt.Errorf("some error"))
			},
			input:   &model.WebUser{Email: "test@test.com", Password: "test_pswd"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateWebUser(tt.input)
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

func Test_GetWebUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewWebUserRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name           string
		mock           func()
		input          string
		want           *model.WebUser
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [web user found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password"}).
					AddRow(1, "test@test.com", "test_pswd")

				mock.ExpectQuery("SELECT * FROM web_user WHERE email = $1;").
					WithArgs("test@test.com").WillReturnRows(rows)
			},
			input: "test@test.com",
			want:  &model.WebUser{ID: 1, Email: "test@test.com", Password: "test_pswd"},
		},
		{
			name: "Error: [web user not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password"})

				mock.ExpectQuery("SELECT * FROM web_user WHERE email = $1;").
					WithArgs("test@test.com").WillReturnRows(rows)
			},
			input:          "test@test.com",
			wantErr:        true,
			expectedErrMsg: "web user not found",
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM web_user WHERE email = $1;").
					WithArgs("test@test.com").WillReturnError(fmt.Errorf("some error"))
			},
			input:          "test@test.com",
			wantErr:        true,
			expectedErrMsg: "failed to get web user by email: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetWebUserByEmail(tt.input)
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
