package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/VladPetriv/scanner_backend_api/internal/handler"
	"github.com/VladPetriv/scanner_backend_api/internal/model"
	"github.com/VladPetriv/scanner_backend_api/internal/service"
	"github.com/VladPetriv/scanner_backend_api/internal/service/mocks"
	"github.com/VladPetriv/scanner_backend_api/internal/store/pg"
	"github.com/VladPetriv/scanner_backend_api/pkg/lib"
	"github.com/VladPetriv/scanner_backend_api/pkg/logger"
)

func Test_GetFullRepliesByMessageIDHandler(t *testing.T) {
	testReplies := []model.FullReplie{
		{
			ID:           1,
			UserID:       1,
			MessageID:    1,
			Title:        "test1",
			UserFullname: "test1 test1",
			UserImageURL: "test1.jpg",
		},
		{
			ID:           2,
			UserID:       2,
			MessageID:    1,
			Title:        "test2",
			UserFullname: "test2 test2",
			UserImageURL: "test2.jpg",
		},
	}

	tests := []struct {
		name            string
		mock            func(replieSrv *mocks.ReplieService)
		input           string
		wantErr         bool
		expectedErr     lib.HttpError
		expectedReplies []model.FullReplie
		expectedCode    int
	}{
		{
			name: "Ok: [replies found]",
			mock: func(replieSrv *mocks.ReplieService) {
				replieSrv.On("GetFullRepliesByMessageID", 1).Return(testReplies, nil)
			},
			input:           "1",
			expectedReplies: testReplies,
			expectedCode:    http.StatusOK,
		},
		{
			name: "Error: [replies not found]",
			mock: func(replieSrv *mocks.ReplieService) {
				replieSrv.On("GetFullRepliesByMessageID", 1).Return(nil, fmt.Errorf("[Replie] srv.GetFullRepliesByMessageID error: %w", pg.ErrFullRepliesNotFound))
			},
			input:        "1",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 404, Name: "Not Found", Message: "full replies not found"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "Error: [some internal error]",
			mock: func(replieSrv *mocks.ReplieService) {
				replieSrv.On("GetFullRepliesByMessageID", 1).Return(nil, fmt.Errorf("[Replie] srv.GetFullRepliesByMessageID error: some error"))
			},
			input:        "1",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 500, Name: "Internal Server Error", Message: "[Replie] srv.GetFullRepliesByMessageID error: some error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "Error: [message id is not valid]",
			mock:         func(replieSrv *mocks.ReplieService) {},
			input:        "hello",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 400, Name: "Bad Request", Message: "message id is not valid"},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/replie/%s", tt.input), nil)
			if err != nil {
				t.Fatalf("could not create request: %s", err)
			}

			rr := httptest.NewRecorder()

			log := logger.Get("debug")

			replieSrv := &mocks.ReplieService{}
			tt.mock(replieSrv)

			handler := handler.New(&service.Manager{Replie: replieSrv}, log)

			router := mux.NewRouter()
			router.HandleFunc("/replie/{message_id}", handler.GetFullRepliesByMessageIDHandler)
			router.ServeHTTP(rr, req)

			decodedReplies := []model.FullReplie{}
			decodedErr := lib.HttpError{}

			if tt.wantErr {
				json.NewDecoder(rr.Body).Decode(&decodedErr)

				assert.EqualValues(t, tt.expectedErr, decodedErr)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			} else {
				json.NewDecoder(rr.Body).Decode(&decodedReplies)

				assert.EqualValues(t, tt.expectedReplies, decodedReplies)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			}
		})
	}
}
