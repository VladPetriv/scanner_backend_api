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

func Test_CreateMessage(t *testing.T) {
	messageInput := &model.MessageDTO{ChannelID: 1, UserID: 1, Title: "test", MessageURL: "test.url", ImageURL: "test.jpg"}

	tests := []struct {
		name           string
		mock           func(messageRepo *mocks.MessageRepo)
		input          *model.MessageDTO
		want           int
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [message created]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("CreateMessage", messageInput).Return(1, nil)
			},
			input: messageInput,
			want:  1,
		},
		{
			name: "Error: [message not created]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("CreateMessage", messageInput).Return(0, pg.ErrMessageNotCreated)
			},
			input:          messageInput,
			wantErr:        true,
			expectedErrMsg: "[Message] srv.CreateMessage error: message not created",
		},
		{
			name: "Error: [some store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("CreateMessage", messageInput).Return(0, fmt.Errorf("failed to create message: some error"))
			},
			input:          messageInput,
			wantErr:        true,
			expectedErrMsg: "[Message] srv.CreateMessage error: failed to create message: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messageRepo := &mocks.MessageRepo{}
			srv := service.NewMessageService(&store.Store{Message: messageRepo})

			tt.mock(messageRepo)

			got, err := srv.CreateMessage(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			messageRepo.AssertExpectations(t)
		})
	}
}

func Test_GetMessagesCount(t *testing.T) {
	tests := []struct {
		name           string
		mock           func(messageRepo *mocks.MessageRepo)
		want           int
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [messages count found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCount").Return(10, nil)
			},
			want: 10,
		},
		{
			name: "Error: [messages count not found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCount").Return(0, fmt.Errorf("failed to get messages count: sql: no rows in result set"))
			},
			wantErr:        true,
			expectedErrMsg: "[Message] srv.GetMessagesCount error: failed to get messages count: sql: no rows in result set",
		},
		{
			name: "Error: [some store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCount").Return(0, fmt.Errorf("failed to get messages count: some error"))
			},
			wantErr:        true,
			expectedErrMsg: "[Message] srv.GetMessagesCount error: failed to get messages count: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messageRepo := &mocks.MessageRepo{}
			srv := service.NewMessageService(&store.Store{Message: messageRepo})

			tt.mock(messageRepo)

			got, err := srv.GetMessagesCount()
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			messageRepo.AssertExpectations(t)
		})
	}
}

func Test_GetMessagesCountByChannelID(t *testing.T) {

	tests := []struct {
		name           string
		mock           func(messageRepo *mocks.MessageRepo)
		input          int
		want           int
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [messages count by channel id found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCountByChannelID", 1).Return(10, nil)
			},
			input: 1,
			want:  10,
		},
		{
			name: "Error: [message count by channel id not found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCountByChannelID", 404).Return(0, fmt.Errorf("failed to get messages count by channel id: sql: no rows in result set"))
			},
			input:          404,
			wantErr:        true,
			expectedErrMsg: "[Message] srv.GetMessagesCountByChannelID error: failed to get messages count by channel id: sql: no rows in result set",
		},
		{
			name: "Error: [some store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetMessagesCountByChannelID", 0).Return(0, fmt.Errorf("failed to get messages count by channel id: some error"))
			},
			input:          0,
			wantErr:        true,
			expectedErrMsg: "[Message] srv.GetMessagesCountByChannelID error: failed to get messages count by channel id: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messageRepo := &mocks.MessageRepo{}
			srv := service.NewMessageService(&store.Store{Message: messageRepo})

			tt.mock(messageRepo)

			got, err := srv.GetMessagesCountByChannelID(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			messageRepo.AssertExpectations(t)
		})
	}
}

func Test_GetFullMessagesByPage(t *testing.T) {
	data := []model.FullMessage{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5},
		{ID: 6}, {ID: 7}, {ID: 8}, {ID: 9}, {ID: 10},
		{ID: 11}, {ID: 12}, {ID: 13}, {ID: 14}, {ID: 15},
		{ID: 16}, {ID: 17}, {ID: 18}, {ID: 19}, {ID: 20},
	}

	tests := []struct {
		name           string
		mock           func(messageRepo *mocks.MessageRepo)
		input          int
		want           []model.FullMessage
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [messages found on page 1]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByPage", 0).Return(data[:10], nil)
			},
			input: 1,
			want:  data[:10],
		},
		{
			name: "Ok: [messages found on page 2]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByPage", 10).Return(data[11:], nil)
			},
			input: 2,
			want:  data[11:],
		},
		{
			name: "Error: [messages not found on page 3]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByPage", 20).Return(nil, fmt.Errorf("full messages not found"))
			},
			input:          3,
			wantErr:        true,
			expectedErrMsg: "[Message] srv.GetFullMessagesByPage error: full messages not found",
		},
		{
			name: "Error: [some store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByPage", 0).Return(nil, fmt.Errorf("failed to get full messages by page: some error"))
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "[Message] srv.GetFullMessagesByPage error: failed to get full messages by page: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messageRepo := &mocks.MessageRepo{}
			srv := service.NewMessageService(&store.Store{Message: messageRepo})

			tt.mock(messageRepo)

			got, err := srv.GetFullMessagesByPage(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			messageRepo.AssertExpectations(t)
		})
	}
}

func Test_GetFullMessagesByChannelIDAndPage(t *testing.T) {
	data := []model.FullMessage{
		{ID: 1, ChannelName: "test1"}, {ID: 2, ChannelName: "test1"}, {ID: 3, ChannelName: "test1"}, {ID: 4, ChannelName: "test1"}, {ID: 5, ChannelName: "test1"},
		{ID: 6, ChannelName: "test1"}, {ID: 7, ChannelName: "test1"}, {ID: 8, ChannelName: "test1"}, {ID: 9, ChannelName: "test1"}, {ID: 10, ChannelName: "test1"},
		{ID: 11, ChannelName: "test2"}, {ID: 12, ChannelName: "test2"}, {ID: 13, ChannelName: "test2"}, {ID: 14, ChannelName: "test2"}, {ID: 15, ChannelName: "test2"},
		{ID: 16, ChannelName: "test2"}, {ID: 17, ChannelName: "test2"}, {ID: 18, ChannelName: "test2"}, {ID: 19, ChannelName: "test2"}, {ID: 20, ChannelName: "test2"},
	}

	tests := []struct {
		name           string
		mock           func(messageRepo *mocks.MessageRepo)
		ID             int
		page           int
		want           []model.FullMessage
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [messages by channel id on page 1 found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByChannelIDAndPage", 1, 0).Return(data[:10], nil)
			},
			ID:   1,
			page: 1,
			want: data[:10],
		},
		{
			name: "Ok: [messages by channel id on page 2 found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByChannelIDAndPage", 2, 10).Return(data[11:], nil)
			},
			ID:   2,
			page: 2,
			want: data[11:],
		},
		{
			name: "Error: [messages by channel id on page 3 not found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByChannelIDAndPage", 1, 20).Return(nil, fmt.Errorf("full messages not found"))
			},
			ID:             1,
			page:           3,
			wantErr:        true,
			expectedErrMsg: "[Message] srv.GetFullMessagesByChannelIDAndPage error: full messages not found",
		},
		{
			name: "Error: [some store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByChannelIDAndPage", 1, 0).Return(nil, fmt.Errorf("failed to get full messages by channel id and page: some error"))
			},
			ID:             1,
			page:           1,
			wantErr:        true,
			expectedErrMsg: "[Message] srv.GetFullMessagesByChannelIDAndPage error: failed to get full messages by channel id and page: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messageRepo := &mocks.MessageRepo{}
			srv := service.NewMessageService(&store.Store{Message: messageRepo})

			tt.mock(messageRepo)

			got, err := srv.GetFullMessagesByChannelIDAndPage(tt.ID, tt.page)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			messageRepo.AssertExpectations(t)
		})
	}
}

func Test_GetFullMessagesByUserID(t *testing.T) {
	data := []model.FullMessage{
		{ID: 1, UserID: 1}, {ID: 2, UserID: 1}, {ID: 3, UserID: 1}, {ID: 4, UserID: 1}, {ID: 5, UserID: 1},
		{ID: 6, UserID: 1}, {ID: 7, UserID: 1}, {ID: 8, UserID: 1}, {ID: 9, UserID: 1}, {ID: 10, UserID: 1},
		{ID: 11, UserID: 2}, {ID: 12, UserID: 2}, {ID: 13, UserID: 2}, {ID: 14, UserID: 2}, {ID: 15, UserID: 2},
		{ID: 16, UserID: 2}, {ID: 17, UserID: 2}, {ID: 18, UserID: 2}, {ID: 19, UserID: 2}, {ID: 20, UserID: 2},
	}

	tests := []struct {
		name           string
		mock           func(messageRepo *mocks.MessageRepo)
		ID             int
		want           []model.FullMessage
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [messages by user id found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByUserID", 1).Return(data[:10], nil)
			},
			ID:   1,
			want: data[:10],
		},
		{
			name: "Ok: [messages by user id found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByUserID", 2).Return(data[11:], nil)
			},
			ID:   2,
			want: data[11:],
		},
		{
			name: "Error: [messages by user id not found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByUserID", 1).Return(nil, fmt.Errorf("full messages not found"))
			},
			ID:             1,
			wantErr:        true,
			expectedErrMsg: "[Message] srv.GetFullMessagesByUserID error: full messages not found",
		},
		{
			name: "Error: [some store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessagesByUserID", 1).Return(nil, fmt.Errorf("failed to get full messages by user ID: some error"))
			},
			ID:             1,
			wantErr:        true,
			expectedErrMsg: "[Message] srv.GetFullMessagesByUserID error: failed to get full messages by user ID: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messageRepo := &mocks.MessageRepo{}
			srv := service.NewMessageService(&store.Store{Message: messageRepo})

			tt.mock(messageRepo)

			got, err := srv.GetFullMessagesByUserID(tt.ID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			messageRepo.AssertExpectations(t)
		})
	}
}

func Test_GetFullMessageByID(t *testing.T) {
	data := &model.FullMessage{
		ID:              1,
		Title:           "test",
		MessageURL:      "test.tg",
		MessageImageURL: "test.jpg",
		ChannelName:     "testc",
		ChannelTitle:    "testc testc",
		ChannelImageURL: "testc.jpg",
		UserID:          1,
		UserFullname:    "testu testu",
		UserImageURL:    "testu.jpg",
		RepliesCount:    0,
	}

	tests := []struct {
		name           string
		mock           func(messageRepo *mocks.MessageRepo)
		input          int
		want           *model.FullMessage
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name: "Ok: [message found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessageByID", 1).Return(data, nil)
			},
			input: 1,
			want:  data,
		},
		{
			name: "Error: [message not found]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessageByID", 404).Return(nil, fmt.Errorf("full messages not found"))
			},
			input:          404,
			wantErr:        true,
			expectedErrMsg: "[Message] srv.GetFullMessageByID error: full messages not found",
		},
		{
			name: "Error: [some store error]",
			mock: func(messageRepo *mocks.MessageRepo) {
				messageRepo.On("GetFullMessageByID", 1).Return(nil, fmt.Errorf("failed to get full message by id: some error"))
			},
			input:          1,
			wantErr:        true,
			expectedErrMsg: "[Message] srv.GetFullMessageByID error: failed to get full message by id: some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messageRepo := &mocks.MessageRepo{}
			srv := service.NewMessageService(&store.Store{Message: messageRepo})

			tt.mock(messageRepo)

			got, err := srv.GetFullMessageByID(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			messageRepo.AssertExpectations(t)
		})
	}
}
