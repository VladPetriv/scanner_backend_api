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

func Test_GetSavedMessages(t *testing.T) {
	data := []model.Saved{
		{ID: 1, UserID: 1, MessageID: 1},
		{ID: 2, UserID: 1, MessageID: 2},
	}

	tests := []struct {
		name           string
		mock           func(savedRepo *mocks.SavedRepo)
		input          int
		want           []model.Saved
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [saved messages found]",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("GetSavedMessages", 1).Return(data, nil)
			},
			input: 1,
			want:  data,
		},
		{
			name: "Error: [saved messages not found]",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("GetSavedMessages", 1).Return(nil, fmt.Errorf("saved messages not found"))
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "[Saved] srv.GetSavedMessages error: saved messages not found",
		},
		{
			name: "Error: [some store error]",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("GetSavedMessages", 1).Return(nil, fmt.Errorf("failed to get saved messages by user id: some error"))
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "[Saved] srv.GetSavedMessages error: failed to get saved messages by user id: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savedRepo := &mocks.SavedRepo{}
			srv := service.NewSavedService(&store.Store{Saved: savedRepo})

			tt.mock(savedRepo)

			got, err := srv.GetSavedMessages(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			savedRepo.AssertExpectations(t)
		})
	}
}

func Test_CreateSavedMessage(t *testing.T) {
	data := &model.Saved{ID: 1, MessageID: 1, UserID: 1}

	tests := []struct {
		name           string
		mock           func(savedRepo *mocks.SavedRepo)
		input          *model.Saved
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [saved message created]",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("CreateSavedMessage", data).Return(1, nil)
			},
			input: data,
		},
		{
			name: "Error: [saved message not created]",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("CreateSavedMessage", data).Return(0, fmt.Errorf("saved message not created"))
			},
			input:          data,
			wantErr:        true,
			expectedErrMsg: "[Saved] srv.CreateSavedMessage error: saved message not created",
		},
		{
			name: "Error: [some store error]",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("CreateSavedMessage", data).Return(0, fmt.Errorf("failed to create saved message: some error"))
			},
			input:          data,
			wantErr:        true,
			expectedErrMsg: "[Saved] srv.CreateSavedMessage error: failed to create saved message: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savedRepo := &mocks.SavedRepo{}
			srv := service.NewSavedService(&store.Store{Saved: savedRepo})

			tt.mock(savedRepo)

			err := srv.CreateSavedMessage(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}

			savedRepo.AssertExpectations(t)
		})
	}
}

func Test_DeleteSavedMessage(t *testing.T) {
	tests := []struct {
		name           string
		mock           func(savedRepo *mocks.SavedRepo)
		input          int
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [saved message deleted]",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("DeleteSavedMessage", 1).Return(1, nil)
			},
			input: 1,
		},
		{
			name: "Error: [saved message not deleted]",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("DeleteSavedMessage", 1).Return(0, fmt.Errorf("saved message not deleted"))
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "[Saved] srv.DeleteSavedMessage error: saved message not deleted",
		},
		{
			name: "Error: [some store error]",
			mock: func(savedRepo *mocks.SavedRepo) {
				savedRepo.On("DeleteSavedMessage", 1).Return(1, fmt.Errorf("failed to delete saved message by id: some error"))
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "[Saved] srv.DeleteSavedMessage error: failed to delete saved message by id: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savedRepo := &mocks.SavedRepo{}
			srv := service.NewSavedService(&store.Store{Saved: savedRepo})

			tt.mock(savedRepo)

			err := srv.DeleteSavedMessage(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}

			savedRepo.AssertExpectations(t)
		})
	}
}
