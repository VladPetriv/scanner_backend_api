package service_test

import (
	"fmt"
	"testing"

	"github.com/VladPetriv/scanner_backend_api/internal/model"
	"github.com/VladPetriv/scanner_backend_api/internal/service"
	"github.com/VladPetriv/scanner_backend_api/internal/store"
	"github.com/VladPetriv/scanner_backend_api/internal/store/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_CreateWebUser(t *testing.T) {
	data := &model.WebUser{
		Email:    "test@test.com",
		Password: "test_pswd",
	}
	tests := []struct {
		name    string
		mock    func(webUserRepo *mocks.WebUserRepo)
		input   *model.WebUser
		wantErr bool
	}{
		{
			name: "Ok: [web user created]",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("CreateWebUser", data).Return(1, nil)
			},
			input: data,
		},
		{
			name: "Error: [web user not created]",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("CreateWebUser", data).Return(0, fmt.Errorf("web user not created"))
			},
			input:   data,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			webUserRepo := &mocks.WebUserRepo{}
			srv := service.NewWebUserService(&store.Store{WebUser: webUserRepo})

			tt.mock(webUserRepo)

			err := srv.CreateWebUser(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			webUserRepo.AssertExpectations(t)
		})
	}
}

func Test_GetWebUserByEmail(t *testing.T) {
	data := &model.WebUser{ID: 1, Email: "test@test.com", Password: "test_pswd"}

	tests := []struct {
		name           string
		mock           func(webUserRepo *mocks.WebUserRepo)
		input          string
		want           *model.WebUser
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [web user found]",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByEmail", "test@test.com").Return(data, nil)
			},
			input: "test@test.com",
			want:  data,
		},
		{
			name: "Error: [web user not found]",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByEmail", "test@test.com").Return(nil, fmt.Errorf("web user not found"))
			},
			input:          "test@test.com",
			wantErr:        true,
			expectedErrMsg: "[WebUser] srv.GetWebUserByEmail error: web user not found",
		},
		{
			name: "Error: [some store error]",
			mock: func(webUserRepo *mocks.WebUserRepo) {
				webUserRepo.On("GetWebUserByEmail", "test@test.com").Return(nil, fmt.Errorf("failed to get web user by email: some error"))
			},
			input:          "test@test.com",
			wantErr:        true,
			expectedErrMsg: "[WebUser] srv.GetWebUserByEmail error: failed to get web user by email: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			webUserRepo := &mocks.WebUserRepo{}
			srv := service.NewWebUserService(&store.Store{WebUser: webUserRepo})

			tt.mock(webUserRepo)

			got, err := srv.GetWebUserByEmail(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			webUserRepo.AssertExpectations(t)
		})
	}
}
