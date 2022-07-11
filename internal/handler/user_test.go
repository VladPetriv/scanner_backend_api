package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VladPetriv/scanner_backend_api/internal/handler"
	"github.com/VladPetriv/scanner_backend_api/internal/model"
	"github.com/VladPetriv/scanner_backend_api/internal/service"
	"github.com/VladPetriv/scanner_backend_api/internal/service/mocks"
	"github.com/VladPetriv/scanner_backend_api/internal/store/pg"
	"github.com/VladPetriv/scanner_backend_api/pkg/lib"
	"github.com/VladPetriv/scanner_backend_api/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_GetUserByIDHandler(t *testing.T) {

	testUser := &model.User{
		ID:       1,
		Username: "test",
		Fullname: "test test",
		ImageURL: "test.jpg",
	}

	tests := []struct {
		name         string
		mock         func(userSrv *mocks.UserService)
		input        string
		wantErr      bool
		expectError  lib.HttpError
		expectedUser model.User
		expectedCode int
	}{
		{
			name: "Ok: [user found]",
			mock: func(userSrv *mocks.UserService) {
				userSrv.On("GetUserByID", 1).Return(testUser, nil)
			},
			input:        "1",
			expectedCode: http.StatusOK,
			expectedUser: *testUser,
		},
		{
			name: "Error: [user not found]",
			mock: func(userSrv *mocks.UserService) {
				userSrv.On("GetUserByID", 1).Return(nil, fmt.Errorf("[User] srv.GetUserByID error: %w", pg.ErrUserNotFound))
			},
			input:        "1",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
			expectError:  lib.HttpError{Code: 404, Name: "Not Found", Message: "user not found"},
		},
		{
			name: "Error: [some internal error]",
			mock: func(userSrv *mocks.UserService) {
				userSrv.On("GetUserByID", 1).Return(nil, fmt.Errorf("[User] srv.GetUserByID error: %w", fmt.Errorf("some internal error")))
			},
			input:        "1",
			wantErr:      true,
			expectedCode: http.StatusInternalServerError,
			expectError:  lib.HttpError{Code: 500, Name: "Internal Server Error", Message: "[User] srv.GetUserByID error: some internal error"},
		},
		{
			name:         "Error: [user id is not valid]",
			mock:         func(userSrv *mocks.UserService) {},
			input:        "hello",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
			expectError:  lib.HttpError{Code: 400, Name: "Bad Request", Message: "user id is not valid"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/user/%s", tt.input), nil)
			if err != nil {
				t.Fatalf("could not create request: %s", err)
			}

			rr := httptest.NewRecorder()

			log := logger.Get("debug")

			userSrv := &mocks.UserService{}
			tt.mock(userSrv)

			handler := handler.New(&service.Manager{User: userSrv}, log)

			router := mux.NewRouter()
			router.HandleFunc("/user/{id}", handler.GetUserByIDHandler)
			router.ServeHTTP(rr, req)

			decodedUser := model.User{}
			decodedErr := lib.HttpError{}

			if tt.wantErr {
				json.NewDecoder(rr.Body).Decode(&decodedErr)

				assert.EqualValues(t, tt.expectError, decodedErr)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			} else {
				json.NewDecoder(rr.Body).Decode(&decodedUser)

				assert.EqualValues(t, tt.expectedUser, decodedUser)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			}
		})
	}
}
