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

func Test_CreateReplie(t *testing.T) {
	replieInput := &model.ReplieDTO{MessageID: 1, UserID: 1, Title: "test", ImageURL: "test.jpg"}

	tests := []struct {
		name           string
		mock           func(replieRepo *mocks.ReplieRepo)
		input          *model.ReplieDTO
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [replie created]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("CreateReplie", replieInput).Return(nil)
			},
			input: replieInput,
		},
		{
			name: "Error: [some store error]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("CreateReplie", replieInput).Return(fmt.Errorf("failed to create replie: some error"))
			},
			input:          replieInput,
			wantErr:        true,
			expectedErrMsg: "[Replie] srv.CreateReplie error: failed to create replie: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			replieRepo := &mocks.ReplieRepo{}
			srv := service.NewReplieService(&store.Store{Replie: replieRepo})

			tt.mock(replieRepo)

			err := srv.CreateReplie(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}

			replieRepo.AssertExpectations(t)
		})
	}
}

func Test_GetFullRepliesByMessageID(t *testing.T) {
	data := []model.FullReplie{
		{ID: 1, Title: "test1", MessageID: 1, UserID: 1, UserFullname: "test1 test1", UserImageURL: "test1.jpg"},
		{ID: 2, Title: "test2", MessageID: 1, UserID: 2, UserFullname: "test2 test2", UserImageURL: "test2.jpg"},
	}

	tests := []struct {
		name           string
		mock           func(replieRepo *mocks.ReplieRepo)
		input          int
		want           []model.FullReplie
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [full replies found]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetFullRepliesByMessageID", 1).Return(data, nil)
			},
			input: 1,
			want:  data,
		},
		{
			name: "Error: [full replies not found]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetFullRepliesByMessageID", 1).Return(nil, fmt.Errorf("full replies not found"))
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "[Replie] srv.GetFullRepliesByMessageID error: full replies not found",
		},
		{
			name: "Error: [some store error]",
			mock: func(replieRepo *mocks.ReplieRepo) {
				replieRepo.On("GetFullRepliesByMessageID", 1).Return(nil, fmt.Errorf("failed to get full replies by message ID: some error"))
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "[Replie] srv.GetFullRepliesByMessageID error: failed to get full replies by message ID: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			replieRepo := &mocks.ReplieRepo{}
			srv := service.NewReplieService(&store.Store{Replie: replieRepo})

			tt.mock(replieRepo)

			got, err := srv.GetFullRepliesByMessageID(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			replieRepo.AssertExpectations(t)
		})
	}
}
