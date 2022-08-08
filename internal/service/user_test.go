package service_test

import (
	"fmt"
	"testing"

	"github.com/VladPetriv/scanner_backend_api/internal/model"
	"github.com/VladPetriv/scanner_backend_api/internal/service"
	"github.com/VladPetriv/scanner_backend_api/internal/store"
	"github.com/VladPetriv/scanner_backend_api/internal/store/mocks"
	"github.com/VladPetriv/scanner_backend_api/internal/store/pg"
	"github.com/stretchr/testify/assert"
)

func Test_CreateUser(t *testing.T) {
	userInput := &model.UserDTO{Username: "test", Fullname: "test test", ImageURL: "test.jpg"}

	tests := []struct {
		name           string
		mock           func(userRepo *mocks.UserRepo)
		input          *model.UserDTO
		want           int
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [new user created]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", userInput.Username).Return(nil, pg.ErrUserNotFound)
				userRepo.On("CreateUser", userInput).Return(1, nil)
			},
			input: userInput,
			want:  1,
		},
		{
			name: "Ok: [new user not created, candidate return an id]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", userInput.Username).Return(&model.User{ID: 10}, nil)
			},
			input: userInput,
			want:  10,
		},
		{
			name: "Error: [user not created]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", userInput.Username).Return(nil, pg.ErrUserNotFound)
				userRepo.On("CreateUser", userInput).Return(0, pg.ErrUserNotCreated)
			},
			input:          userInput,
			wantErr:        true,
			expectedErrMsg: "[User] srv.CreateUser error: user not created",
		},
		{
			name: "Error: [some store error]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByUsername", userInput.Username).Return(nil, pg.ErrUserNotFound)
				userRepo.On("CreateUser", userInput).Return(0, fmt.Errorf("failed to create user: some error"))
			},
			input:          userInput,
			wantErr:        true,
			expectedErrMsg: "[User] srv.CreateUser error: failed to create user: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mocks.UserRepo{}
			srv := service.NewUserService(&store.Store{User: userRepo})

			tt.mock(userRepo)

			got, err := srv.CreateUser(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			userRepo.AssertExpectations(t)
		})
	}
}

func Test_GetUserByID(t *testing.T) {
	data := &model.User{ID: 1, Username: "test", Fullname: "test test", ImageURL: "test.jpg"}

	tests := []struct {
		name           string
		mock           func(userRepo *mocks.UserRepo)
		input          int
		want           *model.User
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [user found]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByID", 1).Return(data, nil)
			},
			input: 1,
			want:  data,
		},
		{
			name: "Error: [user not found]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByID", 1).Return(nil, fmt.Errorf("user not found"))
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "[User] srv.GetUserByID error: user not found",
		},
		{
			name: "Error: [some store error]",
			mock: func(userRepo *mocks.UserRepo) {
				userRepo.On("GetUserByID", 1).Return(nil, fmt.Errorf("failed to get user by id: some error"))
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "[User] srv.GetUserByID error: failed to get user by id: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mocks.UserRepo{}
			srv := service.NewUserService(&store.Store{User: userRepo})

			tt.mock(userRepo)

			got, err := srv.GetUserByID(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			userRepo.AssertExpectations(t)
		})
	}
}
