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

func Test_GetFullRepliesByMessageID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewReplieRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name           string
		mock           func()
		input          int
		want           []model.FullReplie
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [full replies found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_id", "imageurl", "userid", "fullname", "userimageurl"}).
					AddRow(1, "test1", 1, "test1r.jpg", 1, "test1 test1", "test1.jpg").
					AddRow(2, "test2", 1, "test2r.jpg", 2, "test2 test2", "test2.jpg")

				mock.ExpectQuery(
					`SELECT
					r.id, r.title, r.message_id, r.imageurl, 
					u.id as userId, u.fullname, u.imageurl AS userimageurl 	
					FROM replie r 
					LEFT JOIN tg_user u ON u.id = r.user_id
					WHERE r.message_id = $1
					ORDER BY r.id DESC NULLS LAST;`,
				).WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want: []model.FullReplie{
				{ID: 1, Title: "test1", MessageID: 1, ImageURL: "test1r.jpg", UserID: 1, UserFullname: "test1 test1", UserImageURL: "test1.jpg"},
				{ID: 2, Title: "test2", MessageID: 1, ImageURL: "test2r.jpg", UserID: 2, UserFullname: "test2 test2", UserImageURL: "test2.jpg"},
			},
		},
		{
			name: "Error: [full replies not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_id", "imageurl", "userid", "fullname", "userimageurl"})

				mock.ExpectQuery(
					`SELECT
					r.id, r.title, r.message_id, r.imageurl,
					u.id as userId, u.fullname, u.imageurl AS userimageurl	
					FROM replie r 
					LEFT JOIN tg_user u ON u.id = r.user_id
					WHERE r.message_id = $1
					ORDER BY r.id DESC NULLS LAST;`,
				).WithArgs(1).WillReturnRows(rows)
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "full replies not found",
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery(
					`SELECT
					r.id, r.title, r.message_id, r.imageurl,
					u.id as userId, u.fullname, u.imageurl AS userimageurl 	
					FROM replie r 
					LEFT JOIN tg_user u ON u.id = r.user_id
					WHERE r.message_id = $1
					ORDER BY r.id DESC NULLS LAST;`,
				).WithArgs(1).WillReturnError(fmt.Errorf("some error"))
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "failed to get full replies by message ID: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetFullRepliesByMessageID(tt.input)
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
