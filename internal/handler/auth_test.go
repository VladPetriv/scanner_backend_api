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

func Test_SignUpHandler(t *testing.T) {
	testWebUser := model.WebUser{Email: "test@test.com", Password: "test"}

	tests := []struct {
		name           string
		mock           func(userSrv *mocks.WebUserService)
		inputBody      string
		inputUser      model.WebUser
		wantErr        bool
		expectedErr    lib.HttpError
		expectedResult string
		expectedCode   int
	}{
		{
			name: "Ok: [user registered]",
			mock: func(userSrv *mocks.WebUserService) {
				userSrv.On("GetWebUserByEmail", "test@test.com").Return(nil, fmt.Errorf("[WebUser] srv.GetWebUserByEmail error: %w", pg.ErrWebUserNotFound))
				userSrv.On("CreateWebUser", &testWebUser).Return(nil)
			},
			inputUser:      testWebUser,
			inputBody:      `{"email":"test@test.com", "password":"test"}`,
			expectedResult: "\"user created\"\n",
			expectedCode:   http.StatusCreated,
		},
		{
			name: "Error: [user with email is exist]",
			mock: func(userSrv *mocks.WebUserService) {
				userSrv.On("GetWebUserByEmail", "test@test.com").Return(&testWebUser, nil)
			},
			inputUser:    testWebUser,
			inputBody:    `{"email":"test@test.com", "password":"test"}`,
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 409, Name: "Conflict", Message: "user with email test@test.com is exist"},
			expectedCode: http.StatusConflict,
		},
		{
			name: "Error: [some internal error with get by email method]",
			mock: func(userSrv *mocks.WebUserService) {
				userSrv.On("GetWebUserByEmail", "test@test.com").Return(nil, fmt.Errorf("[WebUser] srv.GetWebUserByEmail error: some error"))
			},
			inputUser:    testWebUser,
			inputBody:    `{"email":"test@test.com", "password":"test"}`,
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 500, Name: "Internal Server Error", Message: "[WebUser] srv.GetWebUserByEmail error: some error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "Error: [some internal error with create user method]",
			mock: func(userSrv *mocks.WebUserService) {
				userSrv.On("GetWebUserByEmail", "test@test.com").Return(nil, fmt.Errorf("[WebUser] srv.GetWebUserByEmail error: %w", pg.ErrWebUserNotFound))
				userSrv.On("CreateWebUser", &testWebUser).Return(fmt.Errorf("[WebUser] srv.CreateWebUser error: some error"))
			},
			inputUser:    testWebUser,
			inputBody:    `{"email":"test@test.com", "password":"test"}`,
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 500, Name: "Internal Server Error", Message: "[WebUser] srv.CreateWebUser error: some error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "Error: [request body is not valid]",
			mock:         func(userSrv *mocks.WebUserService) {},
			inputUser:    testWebUser,
			inputBody:    "hello",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 400, Name: "Bad Request", Message: "can't decode request body"},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := bytes.NewBuffer([]byte(tt.inputBody))

			req, err := http.NewRequest(http.MethodPost, "/auth/sign-up", res)
			if err != nil {
				t.Fatalf("could not create request: %s", err)
			}

			rr := httptest.NewRecorder()

			log := logger.Get("debug")

			userSrv := &mocks.WebUserService{}
			tt.mock(userSrv)

			handler := handler.New(&service.Manager{WebUser: userSrv}, log)

			router := mux.NewRouter()
			router.HandleFunc("/auth/sign-up", handler.SignUpHandler)
			router.ServeHTTP(rr, req)

			decodedErr := lib.HttpError{}
			if tt.wantErr {
				json.NewDecoder(rr.Body).Decode(&decodedErr)

				assert.EqualValues(t, tt.expectedErr, decodedErr)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			} else {

				assert.EqualValues(t, tt.expectedResult, rr.Body.String())
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			}

			userSrv.AssertExpectations(t)
		})
	}
}

func Test_SignInHandler(t *testing.T) {
	testWebUser := model.WebUser{Email: "test@test.com", Password: "$2a$14$2nCtguumOv8npwnzjIn6H.refVOG9Fc9AJxUH9Kkkcxtxabhuop2O"}

	tests := []struct {
		name           string
		mock           func(userSrv *mocks.WebUserService, jwtSrv *mocks.JwtService)
		inputBody      string
		wantErr        bool
		expectedErr    lib.HttpError
		expectedResult string
		expectedCode   int
	}{
		{
			name: "Ok: [user is authenticated]",
			mock: func(userSrv *mocks.WebUserService, jwtSrv *mocks.JwtService) {
				userSrv.On("GetWebUserByEmail", "test@test.com").Return(&testWebUser, nil)
				jwtSrv.On("GenerateToken", "test@test.com").Return("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiZW1haWwiOiJ0ZXN0QHRlc3QuY29tIiwiaWF0IjoxNTE2MjM5MDIyfQ.j7o5o8GBkybaYXdFJIi8O6mPF50E-gJWZ3reLfMQD68", nil)
			},
			inputBody:      `{"email":"test@test.com", "password":"test"}`,
			expectedResult: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiZW1haWwiOiJ0ZXN0QHRlc3QuY29tIiwiaWF0IjoxNTE2MjM5MDIyfQ.j7o5o8GBkybaYXdFJIi8O6mPF50E-gJWZ3reLfMQD68",
			expectedCode:   http.StatusOK,
		},
		{
			name: "Error: [uesr not found]",
			mock: func(userSrv *mocks.WebUserService, jwtSrv *mocks.JwtService) {
				userSrv.On("GetWebUserByEmail", "test@test.com").Return(nil, fmt.Errorf("[WebUser] srv.GetWebUserByEmail error: %w", pg.ErrWebUserNotFound))
			},
			inputBody:    `{"email":"test@test.com", "password":"test"}`,
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 404, Name: "Not Found", Message: "user not found"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "Error: [some internal error with get by email method]",
			mock: func(userSrv *mocks.WebUserService, jwtSrv *mocks.JwtService) {
				userSrv.On("GetWebUserByEmail", "test@test.com").Return(nil, fmt.Errorf("[WebUser] srv.GetWebUserByEmail error: some error"))
			},
			inputBody:    `{"email":"test@test.com", "password":"test"}`,
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 500, Name: "Internal Server Error", Message: "[WebUser] srv.GetWebUserByEmail error: some error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "Error: [password is incorrect]",
			mock: func(userSrv *mocks.WebUserService, jwtSrv *mocks.JwtService) {
				userSrv.On("GetWebUserByEmail", "test@test.com").Return(&testWebUser, nil)
			},
			inputBody:    `{"email":"test@test.com", "password":"testx"}`,
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 401, Name: "Unauthorized", Message: "password is incorrect"},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "Error: [some internal error with jwt method]",
			mock: func(userSrv *mocks.WebUserService, jwtSrv *mocks.JwtService) {
				userSrv.On("GetWebUserByEmail", "test@test.com").Return(&testWebUser, nil)
				jwtSrv.On("GenerateToken", "test@test.com").Return("", fmt.Errorf("some error"))
			},
			inputBody:    `{"email":"test@test.com", "password":"test"}`,
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 500, Name: "Internal Server Error", Message: "failed to generate jwt token"},
			expectedCode: 500,
		},
		{
			name:         "Error: [request body is not valid]",
			mock:         func(userSrv *mocks.WebUserService, jwtSrv *mocks.JwtService) {},
			inputBody:    "hello",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 400, Name: "Bad Request", Message: "can't decode request body"},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := bytes.NewBuffer([]byte(tt.inputBody))

			req, err := http.NewRequest(http.MethodPost, "/auth/sign-in", res)
			if err != nil {
				t.Fatalf("could not create request: %s", err)
			}

			rr := httptest.NewRecorder()

			log := logger.Get("debug")

			userSrv := &mocks.WebUserService{}
			jwtSrv := &mocks.JwtService{}
			tt.mock(userSrv, jwtSrv)

			handler := handler.New(&service.Manager{WebUser: userSrv, Jwt: jwtSrv}, log)

			router := mux.NewRouter()
			router.HandleFunc("/auth/sign-in", handler.SignInHandler)
			router.ServeHTTP(rr, req)

			decordedData := map[string]string{}
			decodedErr := lib.HttpError{}

			if tt.wantErr {
				json.NewDecoder(rr.Body).Decode(&decodedErr)

				assert.EqualValues(t, tt.expectedErr, decodedErr)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			} else {
				json.NewDecoder(rr.Body).Decode(&decordedData)

				assert.EqualValues(t, tt.expectedResult, decordedData["token"])
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			}

			userSrv.AssertExpectations(t)
			jwtSrv.AssertExpectations(t)
		})
	}
}
