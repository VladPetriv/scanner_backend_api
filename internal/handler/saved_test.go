package handler_test

import (
	"bytes"
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

func Test_GetSavedMessagesHandler(t *testing.T) {
	testMessages := []model.Saved{
		{ID: 1, UserID: 1, MessageID: 1},
		{ID: 2, UserID: 1, MessageID: 2},
		{ID: 3, UserID: 1, MessageID: 3},
	}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTc5MDgyMTgsImlhdCI6MTY1Nzg2NTAxOCwiVXNlckVtYWlsIjoidGVzdEB0ZXN0LmNvbSJ9.O-CpA-vp3mOuBKybKZBPeIlebTozyxx1_ql8F3P1YzI"

	tests := []struct {
		name           string
		mock           func(savedSrv *mocks.SavedService, jwtSrv *mocks.JwtService, token string)
		token          string
		input          string
		wantErr        bool
		expectedErr    lib.HttpError
		expectedResult []model.Saved
		expectedCode   int
	}{
		{
			name: "Ok: [saved message's found]",
			mock: func(savedSrv *mocks.SavedService, jwtSrv *mocks.JwtService, token string) {
				jwtSrv.On("ParseToken", token).Return("test@test.com", nil)
				savedSrv.On("GetSavedMessages", 1).Return(testMessages, nil)
			},
			input:          "1",
			token:          token,
			expectedResult: testMessages,
			expectedCode:   http.StatusOK,
		},
		{
			name: "Error: [saved message's not found]",
			mock: func(savedSrv *mocks.SavedService, jwtSrv *mocks.JwtService, token string) {
				jwtSrv.On("ParseToken", token).Return("test@test.com", nil)
				savedSrv.On("GetSavedMessages", 1).Return(nil, fmt.Errorf("[Saved] srv.GetSavedMessages error: %w", pg.ErrSavedMessagesNotFound))
			},
			input:        "1",
			token:        token,
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 404, Name: "Not Found", Message: "saved messages not found"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "Error: [some internal error]",
			mock: func(savedSrv *mocks.SavedService, jwtSrv *mocks.JwtService, token string) {
				jwtSrv.On("ParseToken", token).Return("test@test.com", nil)
				savedSrv.On("GetSavedMessages", 1).Return(nil, fmt.Errorf("[Saved] srv.GetSavedMessages error: some error"))
			},
			input:        "1",
			token:        token,
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 500, Name: "Internal Server Error", Message: "[Saved] srv.GetSavedMessages error: some error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "Error: [user id is not valid]",
			mock: func(savedSrv *mocks.SavedService, jwtSrv *mocks.JwtService, token string) {
				jwtSrv.On("ParseToken", token).Return("test@test.com", nil)
			},
			input:        "hello",
			token:        token,
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 400, Name: "Bad Request", Message: "user id is not valid"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Error: [user is not authorized]",
			mock:         func(savedSrv *mocks.SavedService, jwtSrv *mocks.JwtService, token string) {},
			input:        "1",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 401, Name: "Unauthorized", Message: "user is not authorized"},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/saved/%s", tt.input), nil)
			if err != nil {
				t.Fatalf("could not create request: %s", err)
			}

			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))

			}

			rr := httptest.NewRecorder()

			log := logger.Get("debug")

			savedSrv := &mocks.SavedService{}
			jwtSrv := &mocks.JwtService{}
			tt.mock(savedSrv, jwtSrv, tt.token)

			handler := handler.New(&service.Manager{Saved: savedSrv, Jwt: jwtSrv}, log)

			router := mux.NewRouter()
			router.Use(handler.AuthenticateMiddleware)
			router.HandleFunc("/saved/{user_id}", handler.GetSavedMessagesHandler)
			router.ServeHTTP(rr, req)

			decodedResult := []model.Saved{}
			decodedErr := lib.HttpError{}

			if tt.wantErr {
				json.NewDecoder(rr.Body).Decode(&decodedErr)

				assert.EqualValues(t, tt.expectedErr, decodedErr)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			} else {
				json.NewDecoder(rr.Body).Decode(&decodedResult)

				assert.EqualValues(t, tt.expectedResult, decodedResult)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
				assert.EqualValues(t, "test@test.com", rr.Header().Get("email"))
			}

			savedSrv.AssertExpectations(t)
			jwtSrv.AssertExpectations(t)
		})
	}
}

func Test_CreateSavedMessagesHandler(t *testing.T) {
	testMessage := model.Saved{UserID: 1, MessageID: 1}

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTc5MDgyMTgsImlhdCI6MTY1Nzg2NTAxOCwiVXNlckVtYWlsIjoidGVzdEB0ZXN0LmNvbSJ9.O-CpA-vp3mOuBKybKZBPeIlebTozyxx1_ql8F3P1YzI"

	tests := []struct {
		name           string
		mock           func(savedSrv *mocks.SavedService, jwtSrv *mocks.JwtService, token string)
		token          string
		input          string
		wantErr        bool
		expectedErr    lib.HttpError
		expectedResult string
		expectedCode   int
	}{
		{
			name: "Ok: [saved message created]",
			mock: func(savedSrv *mocks.SavedService, jwtSrv *mocks.JwtService, token string) {
				jwtSrv.On("ParseToken", token).Return("test@test.com", nil)
				savedSrv.On("CreateSavedMessage", &testMessage).Return(nil)
			},
			input:          `{"userId":1, "messageId":1}`,
			token:          token,
			expectedResult: "\"saved message created\"\n",
			expectedCode:   http.StatusCreated,
		},
		{
			name: "Error: [some internal error]",
			mock: func(savedSrv *mocks.SavedService, jwtSrv *mocks.JwtService, token string) {
				jwtSrv.On("ParseToken", token).Return("test@test.com", nil)
				savedSrv.On("CreateSavedMessage", &testMessage).Return(fmt.Errorf("[Saved] srv.CreateSavedMessage error: some error"))
			},
			input:        `{"userId":1, "messageId":1}`,
			token:        token,
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 500, Name: "Internal Server Error", Message: "[Saved] srv.CreateSavedMessage error: some error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "Error: [request body is not valid]",
			mock: func(savedSrv *mocks.SavedService, jwtSrv *mocks.JwtService, token string) {
				jwtSrv.On("ParseToken", token).Return("test@test.com", nil)
			},
			input:        "hello",
			token:        token,
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 400, Name: "Bad Request", Message: "can't decode request body"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Error: [user is unauthorized]",
			mock:         func(savedSrv *mocks.SavedService, jwtSrv *mocks.JwtService, token string) {},
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 401, Name: "Unauthorized", Message: "user is not authorized"},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewBuffer([]byte(tt.input))

			req, err := http.NewRequest(http.MethodPost, "/saved/create", body)
			if err != nil {
				t.Fatalf("could not create request: %s", err)
			}

			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))

			}

			rr := httptest.NewRecorder()

			log := logger.Get("debug")

			savedSrv := &mocks.SavedService{}
			jwtSrv := &mocks.JwtService{}
			tt.mock(savedSrv, jwtSrv, tt.token)

			handler := handler.New(&service.Manager{Saved: savedSrv, Jwt: jwtSrv}, log)

			router := mux.NewRouter()
			router.Use(handler.AuthenticateMiddleware)
			router.HandleFunc("/saved/create", handler.CreateSavedMessageHandler)
			router.ServeHTTP(rr, req)

			decodedErr := lib.HttpError{}

			if tt.wantErr {
				json.NewDecoder(rr.Body).Decode(&decodedErr)

				assert.EqualValues(t, tt.expectedErr, decodedErr)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			} else {
				assert.EqualValues(t, tt.expectedResult, rr.Body.String())
				assert.EqualValues(t, tt.expectedCode, rr.Code)
				assert.EqualValues(t, "test@test.com", rr.Header().Get("email"))
			}

			savedSrv.AssertExpectations(t)
			jwtSrv.AssertExpectations(t)
		})
	}
}

func Test_DeleteSavedMessage(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTc5MDgyMTgsImlhdCI6MTY1Nzg2NTAxOCwiVXNlckVtYWlsIjoidGVzdEB0ZXN0LmNvbSJ9.O-CpA-vp3mOuBKybKZBPeIlebTozyxx1_ql8F3P1YzI"

	tests := []struct {
		name           string
		mock           func(savedSrv *mocks.SavedService, jwtSrv *mocks.JwtService, token string)
		token          string
		input          string
		wantErr        bool
		expectedErr    lib.HttpError
		expectedResult string
		expectedCode   int
	}{
		{
			name: "Ok: [saved message created]",
			mock: func(savedSrv *mocks.SavedService, jwtSrv *mocks.JwtService, token string) {
				jwtSrv.On("ParseToken", token).Return("test@test.com", nil)
				savedSrv.On("DeleteSavedMessage", 1).Return(nil)
			},
			input:          "1",
			token:          token,
			expectedResult: "\"saved message deleted\"\n",
			expectedCode:   http.StatusOK,
		},
		{
			name: "Error: [some internal error]",
			mock: func(savedSrv *mocks.SavedService, jwtSrv *mocks.JwtService, token string) {
				jwtSrv.On("ParseToken", token).Return("test@test.com", nil)
				savedSrv.On("DeleteSavedMessage", 1).Return(fmt.Errorf("[Saved] srv.DeleteSavedMessage error: some error"))
			},
			input:        "1",
			token:        token,
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 500, Name: "Internal Server Error", Message: "[Saved] srv.DeleteSavedMessage error: some error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "Error: [message id is not valid]",
			mock: func(savedSrv *mocks.SavedService, jwtSrv *mocks.JwtService, token string) {
				jwtSrv.On("ParseToken", token).Return("test@test.com", nil)
			},
			input:        "hello",
			token:        token,
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 400, Name: "Bad Request", Message: "message id is not valid"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Error: [user is unauthorized]",
			mock:         func(savedSrv *mocks.SavedService, jwtSrv *mocks.JwtService, token string) {},
			input:        "1",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 401, Name: "Unauthorized", Message: "user is not authorized"},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/saved/delete/%s", tt.input), nil)
			if err != nil {
				t.Fatalf("could not create request: %s", err)
			}

			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))

			}

			rr := httptest.NewRecorder()

			log := logger.Get("debug")

			savedSrv := &mocks.SavedService{}
			jwtSrv := &mocks.JwtService{}
			tt.mock(savedSrv, jwtSrv, tt.token)

			handler := handler.New(&service.Manager{Saved: savedSrv, Jwt: jwtSrv}, log)

			router := mux.NewRouter()
			router.Use(handler.AuthenticateMiddleware)
			router.HandleFunc("/saved/delete/{message_id}", handler.DeleteSavedMessageHandler)
			router.ServeHTTP(rr, req)

			decodedErr := lib.HttpError{}

			if tt.wantErr {
				json.NewDecoder(rr.Body).Decode(&decodedErr)

				assert.EqualValues(t, tt.expectedErr, decodedErr)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			} else {
				assert.EqualValues(t, tt.expectedResult, rr.Body.String())
				assert.EqualValues(t, tt.expectedCode, rr.Code)
				assert.EqualValues(t, "test@test.com", rr.Header().Get("email"))
			}

			savedSrv.AssertExpectations(t)
			jwtSrv.AssertExpectations(t)
		})
	}
}
