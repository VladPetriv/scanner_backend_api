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

func Test_GetChannelsCount(t *testing.T) {
	tests := []struct {
		name           string
		mock           func(channelRepo *mocks.ChannelRepo)
		want           int
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [count of channels found]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelsCount").Return(10, nil)
			},
			want: 10,
		},
		{
			name: "Error: [some store error]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelsCount").Return(0, fmt.Errorf("failed to get count of channels: some error"))
			},
			wantErr:        true,
			expectedErrMsg: "[Channel] srv.GetChannelsCount error: failed to get count of channels: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			channelRepo := &mocks.ChannelRepo{}
			srv := service.NewChannelService(&store.Store{Channel: channelRepo})

			tt.mock(channelRepo)

			got, err := srv.GetChannelsCount()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			channelRepo.AssertExpectations(t)
		})
	}
}

func Test_GetChannelsByPage(t *testing.T) {
	data := []model.Channel{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5},
		{ID: 6}, {ID: 7}, {ID: 8}, {ID: 9}, {ID: 10},
		{ID: 11}, {ID: 12}, {ID: 13}, {ID: 14}, {ID: 15},
		{ID: 16}, {ID: 17}, {ID: 18}, {ID: 19}, {ID: 20},
	}

	tests := []struct {
		name           string
		mock           func(channelRepo *mocks.ChannelRepo)
		input          int
		want           []model.Channel
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [channels on page 1 found with input 0]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelsByPage", 0).Return(data[:10], nil)
			},
			input: 0,
			want:  data[:10],
		},
		{
			name: "Ok: [channels on page 1 found with input 1]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelsByPage", 0).Return(data[:10], nil)
			},
			input: 1,
			want:  data[:10],
		},
		{
			name: "Ok: [channels on page 2 found with input 2]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelsByPage", 10).Return(data[11:], nil)
			},
			input: 2,
			want:  data[11:],
		},
		{

			name: "Error: [channels on page 3 not found with input ]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelsByPage", 20).Return(nil, fmt.Errorf("channel/s not found"))
			},
			input:          3,
			wantErr:        true,
			expectedErrMsg: "[Channel] srv.GetChannelsByPage error: channel/s not found",
		},
		{
			name: "Error: [some store error]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelsByPage", 0).Return(nil, fmt.Errorf("failed to get channels by page: some error"))
			},
			input:          0,
			wantErr:        true,
			expectedErrMsg: "[Channel] srv.GetChannelsByPage error: failed to get channels by page: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			channelRepo := &mocks.ChannelRepo{}
			srv := service.NewChannelService(&store.Store{Channel: channelRepo})

			tt.mock(channelRepo)

			got, err := srv.GetChannelsByPage(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			channelRepo.AssertExpectations(t)
		})
	}
}

func Test_GetChannelByName(t *testing.T) {
	data := &model.Channel{ID: 1, Name: "test", Title: "test", ImageURL: "test.jpg"}

	tests := []struct {
		name           string
		mock           func(channelRepo *mocks.ChannelRepo)
		input          string
		want           *model.Channel
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [channel found]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", "test").Return(data, nil)
			},
			input: "test",
			want:  data,
		},
		{
			name: "Error: [channele not found]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", "test").Return(nil, fmt.Errorf("failed to get channel by name: sql: no rows in result set"))
			},
			input:          "test",
			wantErr:        true,
			expectedErrMsg: "[Channel] srv.GetChannelByName error: failed to get channel by name: sql: no rows in result set",
		},
		{
			name: "Error: [some store error]",
			mock: func(channelRepo *mocks.ChannelRepo) {
				channelRepo.On("GetChannelByName", "test").Return(nil, fmt.Errorf("failed to get channel by name: some error"))
			},
			input:          "test",
			wantErr:        true,
			expectedErrMsg: "[Channel] srv.GetChannelByName error: failed to get channel by name: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			channelRepo := &mocks.ChannelRepo{}
			srv := service.NewChannelService(&store.Store{Channel: channelRepo})

			tt.mock(channelRepo)

			got, err := srv.GetChannelByName(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}
