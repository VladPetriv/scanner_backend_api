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

func Test_GetMessagesCountHandler(t *testing.T) {
	tests := []struct {
		name          string
		mock          func(messageSrv *mocks.MessageService)
		wantErr       bool
		expectedErr   lib.HttpError
		expectedCount int
		expectedCode  int
	}{
		{
			name: "Ok: [messages count found]",
			mock: func(messageSrv *mocks.MessageService) {
				messageSrv.On("GetMessagesCount").Return(10, nil)
			},
			expectedCount: 10,
			expectedCode:  http.StatusOK,
		},
		{
			name: "Error: [messages count not found]",
			mock: func(messageSrv *mocks.MessageService) {
				messageSrv.On("GetMessagesCount").Return(0, fmt.Errorf("[Message] srv.GetMessagesCount error: %w", pg.ErrMessagesCountNotFound))
			},
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 404, Name: "Not Found", Message: "messages count not found"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "Error: [some internal error]",
			mock: func(messageSrv *mocks.MessageService) {
				messageSrv.On("GetMessagesCount").Return(0, fmt.Errorf("[Message] srv.GetMessagesCount error: some error"))
			},
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 500, Name: "Internal Server Error", Message: "[Message] srv.GetMessagesCount error: some error"},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/message/count", nil)
			if err != nil {
				t.Fatalf("could not create request: %s", err)
			}

			rr := httptest.NewRecorder()

			log := logger.Get("debug")

			messageSrv := &mocks.MessageService{}
			tt.mock(messageSrv)

			handler := handler.New(&service.Manager{Message: messageSrv}, log)

			router := mux.NewRouter()
			router.HandleFunc("/message/count", handler.GetMessagesCountHandler)
			router.ServeHTTP(rr, req)

			decodedCount := map[string]int{}
			decodedErr := lib.HttpError{}

			if tt.wantErr {
				json.NewDecoder(rr.Body).Decode(&decodedErr)

				assert.EqualValues(t, tt.expectedErr, decodedErr)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			} else {
				json.NewDecoder(rr.Body).Decode(&decodedCount)

				assert.EqualValues(t, tt.expectedCount, decodedCount["count"])
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			}

			messageSrv.AssertExpectations(t)
		})
	}
}

func Test_GetMessagesCountByChannelIDHandler(t *testing.T) {
	tests := []struct {
		name          string
		mock          func(messageSrv *mocks.MessageService)
		input         string
		wantErr       bool
		expectedErr   lib.HttpError
		expectedCount int
		expectedCode  int
	}{
		{
			name: "Ok: [messages count found]",
			mock: func(messageSrv *mocks.MessageService) {
				messageSrv.On("GetMessagesCountByChannelID", 1).Return(10, nil)
			},
			input:         "1",
			expectedCount: 10,
			expectedCode:  http.StatusOK,
		},
		{
			name: "Error: [messages count not found]",
			mock: func(messageSrv *mocks.MessageService) {
				messageSrv.On("GetMessagesCountByChannelID", 1).Return(0, fmt.Errorf("[Message] srv.GetMessagesCountByChannelID error: %w", pg.ErrMessagesCountNotFound))
			},
			input:        "1",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 404, Name: "Not Found", Message: "messages count not found"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "Error: [some internal error]",
			mock: func(messageSrv *mocks.MessageService) {
				messageSrv.On("GetMessagesCountByChannelID", 1).Return(0, fmt.Errorf("[Message] srv.GetMessagesCountByChannelID error: some error"))
			},
			input:        "1",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 500, Name: "Internal Server Error", Message: "[Message] srv.GetMessagesCountByChannelID error: some error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "Error: [channel id is not valid]",
			mock:         func(messageSrv *mocks.MessageService) {},
			input:        "hello",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 400, Name: "Bad Request", Message: "channel id is not valid"},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/message/count/%s", tt.input), nil)
			if err != nil {
				t.Fatalf("could not create request: %s", err)
			}

			rr := httptest.NewRecorder()

			log := logger.Get("debug")

			messageSrv := &mocks.MessageService{}
			tt.mock(messageSrv)

			handler := handler.New(&service.Manager{Message: messageSrv}, log)

			router := mux.NewRouter()
			router.HandleFunc("/message/count/{channel_id}", handler.GetMessagesCountByChannelIDHandler)
			router.ServeHTTP(rr, req)

			decodedCount := map[string]int{}
			decodedErr := lib.HttpError{}

			if tt.wantErr {
				json.NewDecoder(rr.Body).Decode(&decodedErr)

				assert.EqualValues(t, tt.expectedErr, decodedErr)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			} else {
				json.NewDecoder(rr.Body).Decode(&decodedCount)

				assert.EqualValues(t, tt.expectedCount, decodedCount["count"])
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			}

			messageSrv.AssertExpectations(t)
		})
	}

}

func Test_GetFullMessagesByPageHandler(t *testing.T) {
	testMessages := []model.FullMessage{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5},
	}

	tests := []struct {
		name             string
		mock             func(messageSrv *mocks.MessageService)
		input            string
		wantErr          bool
		expectedErr      lib.HttpError
		expectedMessages []model.FullMessage
		expectedCode     int
	}{
		{
			name: "Ok: [messages found]",
			mock: func(messageSrv *mocks.MessageService) {
				messageSrv.On("GetFullMessagesByPage", 1).Return(testMessages, nil)
			},
			input:            "1",
			expectedMessages: testMessages,
			expectedCode:     http.StatusOK,
		},
		{
			name: "Error: [messages not found]",
			mock: func(messageSrv *mocks.MessageService) {
				messageSrv.On("GetFullMessagesByPage", 1).Return(nil, fmt.Errorf("[Message] srv.GetFullMessagesByPage error: %w", pg.ErrFullMessagesNotFound))
			},
			input:        "1",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 404, Name: "Not Found", Message: "full messages not found"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "Error: [some internal error]",
			mock: func(messageSrv *mocks.MessageService) {
				messageSrv.On("GetFullMessagesByPage", 1).Return(nil, fmt.Errorf("[Message] srv.GetFullMessagesByPage error: some error"))
			},
			input:        "1",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 500, Name: "Internal Server Error", Message: "[Message] srv.GetFullMessagesByPage error: some error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "Error: [page is not valid]",
			mock:         func(messageSrv *mocks.MessageService) {},
			input:        "hello",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 400, Name: "Bad Request", Message: "page is not valid"},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/message/?page=%s", tt.input), nil)
			if err != nil {
				t.Fatalf("could not create request: %s", err)
			}

			rr := httptest.NewRecorder()

			log := logger.Get("debug")

			messageSrv := &mocks.MessageService{}
			tt.mock(messageSrv)

			handler := handler.New(&service.Manager{Message: messageSrv}, log)

			router := mux.NewRouter()
			router.HandleFunc("/message/", handler.GetFullMessagesByPageHandler)
			router.ServeHTTP(rr, req)

			decodedMessages := []model.FullMessage{}
			decodedErr := lib.HttpError{}

			if tt.wantErr {
				json.NewDecoder(rr.Body).Decode(&decodedErr)

				assert.EqualValues(t, tt.expectedErr, decodedErr)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			} else {
				json.NewDecoder(rr.Body).Decode(&decodedMessages)

				assert.EqualValues(t, tt.expectedMessages, decodedMessages)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			}

			messageSrv.AssertExpectations(t)
		})
	}
}

func Test_GetFullMessagesByChannelIDAndPageHandler(t *testing.T) {
	testMessages := []model.FullMessage{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5},
		{ID: 6}, {ID: 7}, {ID: 8}, {ID: 9}, {ID: 10},
		{ID: 11}, {ID: 12}, {ID: 13}, {ID: 14}, {ID: 15},
		{ID: 16}, {ID: 17}, {ID: 18}, {ID: 19}, {ID: 20},
	}

	tests := []struct {
		name             string
		mock             func(messageSrv *mocks.MessageService)
		page             string
		channelID        string
		wantErr          bool
		expectedErr      lib.HttpError
		expectedMessages []model.FullMessage
		expectedCode     int
	}{
		{
			name: "Ok: [messages found on page 1]",
			mock: func(messageSrv *mocks.MessageService) {
				messageSrv.On("GetFullMessagesByChannelIDAndPage", 1, 1).Return(testMessages[:10], nil)
			},
			page:             "1",
			channelID:        "1",
			expectedMessages: testMessages[:10],
			expectedCode:     http.StatusOK,
		},
		{
			name: "Ok: [messages found on page 2]",
			mock: func(messageSrv *mocks.MessageService) {
				messageSrv.On("GetFullMessagesByChannelIDAndPage", 1, 2).Return(testMessages[10:], nil)
			},
			page:             "2",
			channelID:        "1",
			expectedMessages: testMessages[10:],
			expectedCode:     http.StatusOK,
		},
		{
			name: "Error: [messages not found]",
			mock: func(messageSrv *mocks.MessageService) {
				messageSrv.On("GetFullMessagesByChannelIDAndPage", 1, 1).Return(nil, fmt.Errorf("[Message] srv.GetFullMessagesByChannelIDAndPage error: %w", pg.ErrFullMessagesNotFound))
			},
			page:         "1",
			channelID:    "1",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 404, Name: "Not Found", Message: "full messages not found"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "Error: [some internal error]",
			mock: func(messageSrv *mocks.MessageService) {
				messageSrv.On("GetFullMessagesByChannelIDAndPage", 1, 1).Return(nil, fmt.Errorf("[Message] srv.GetFullMessagesByChannelIDAndPage error: some error"))
			},
			page:         "1",
			channelID:    "1",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 500, Name: "Internal Server Error", Message: "[Message] srv.GetFullMessagesByChannelIDAndPage error: some error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "Error: [page is not valid]",
			mock:         func(messageSrv *mocks.MessageService) {},
			page:         "hello",
			channelID:    "1",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 400, Name: "Bad Request", Message: "page is not valid"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Error: [channel id is not valid]",
			mock:         func(messageSrv *mocks.MessageService) {},
			page:         "1",
			channelID:    "hello",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 400, Name: "Bad Request", Message: "channel id is not valid"},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/message/channel/%s?page=%s", tt.channelID, tt.page), nil)
			if err != nil {
				t.Fatalf("could not create request: %s", err)
			}

			rr := httptest.NewRecorder()

			log := logger.Get("debug")

			messageSrv := &mocks.MessageService{}
			tt.mock(messageSrv)

			handler := handler.New(&service.Manager{Message: messageSrv}, log)

			router := mux.NewRouter()
			router.HandleFunc("/message/channel/{channel_id}", handler.GetFullMessagesByChannelIDAndPageHandler)
			router.ServeHTTP(rr, req)

			decodedMessages := []model.FullMessage{}
			decodedErr := lib.HttpError{}

			if tt.wantErr {
				json.NewDecoder(rr.Body).Decode(&decodedErr)

				assert.EqualValues(t, tt.expectedErr, decodedErr)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			} else {
				json.NewDecoder(rr.Body).Decode(&decodedMessages)

				assert.EqualValues(t, tt.expectedMessages, decodedMessages)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			}

			messageSrv.AssertExpectations(t)
		})
	}
}

func Test_GetFullMessagesByUserIDHandler(t *testing.T) {
	testMessages := []model.FullMessage{
		{ID: 1, UserID: 1}, {ID: 2, UserID: 1}, {ID: 3, UserID: 1}, {ID: 4, UserID: 1}, {ID: 5, UserID: 1},
		{ID: 6, UserID: 1}, {ID: 7, UserID: 1}, {ID: 8, UserID: 1}, {ID: 9, UserID: 1}, {ID: 10, UserID: 1},
		{ID: 11, UserID: 2}, {ID: 12, UserID: 2}, {ID: 13, UserID: 2}, {ID: 14, UserID: 2}, {ID: 15, UserID: 2},
		{ID: 16, UserID: 2}, {ID: 17, UserID: 2}, {ID: 18, UserID: 2}, {ID: 19, UserID: 2}, {ID: 20, UserID: 2},
	}

	tests := []struct {
		name             string
		mock             func(messageSrv *mocks.MessageService)
		input            string
		wantErr          bool
		expectedErr      lib.HttpError
		expectedMessages []model.FullMessage
		expectedCode     int
	}{
		{
			name: "Ok: [messages found user id = 1]",
			mock: func(messageSrv *mocks.MessageService) {
				messageSrv.On("GetFullMessagesByUserID", 1).Return(testMessages[:10], nil)
			},
			input:            "1",
			expectedMessages: testMessages[:10],
			expectedCode:     http.StatusOK,
		},
		{
			name: "Ok: [messages found user id = 2]",
			mock: func(messageSrv *mocks.MessageService) {
				messageSrv.On("GetFullMessagesByUserID", 2).Return(testMessages[10:], nil)
			},
			input:            "2",
			expectedMessages: testMessages[10:],
			expectedCode:     http.StatusOK,
		},
		{
			name: "Error: [messages not found]",
			mock: func(messageSrv *mocks.MessageService) {
				messageSrv.On("GetFullMessagesByUserID", 1).Return(nil, fmt.Errorf("[Message] srv.GetFullMessagesByUserID error: %w", pg.ErrFullMessagesNotFound))
			},
			input:        "1",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 404, Name: "Not Found", Message: "full messages not found"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "Error: [some internal error]",
			mock: func(messageSrv *mocks.MessageService) {
				messageSrv.On("GetFullMessagesByUserID", 1).Return(nil, fmt.Errorf("[Message] srv.GetFullMessagesByUserID error: some error"))
			},
			input:        "1",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 500, Name: "Internal Server Error", Message: "[Message] srv.GetFullMessagesByUserID error: some error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "Error: [user id is not valid]",
			mock:         func(messageSrv *mocks.MessageService) {},
			input:        "hello",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 400, Name: "Bad Request", Message: "user id is not valid"},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/message/user/%s", tt.input), nil)
			if err != nil {
				t.Fatalf("could not create request: %s", err)
			}

			rr := httptest.NewRecorder()

			log := logger.Get("debug")

			messageSrv := &mocks.MessageService{}
			tt.mock(messageSrv)

			handler := handler.New(&service.Manager{Message: messageSrv}, log)

			router := mux.NewRouter()
			router.HandleFunc("/message/user/{user_id}", handler.GetFullMessagesByUserIDHandler)
			router.ServeHTTP(rr, req)

			decodedMessages := []model.FullMessage{}
			decodedErr := lib.HttpError{}

			if tt.wantErr {
				json.NewDecoder(rr.Body).Decode(&decodedErr)

				assert.EqualValues(t, tt.expectedErr, decodedErr)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			} else {
				json.NewDecoder(rr.Body).Decode(&decodedMessages)

				assert.EqualValues(t, tt.expectedMessages, decodedMessages)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			}

			messageSrv.AssertExpectations(t)
		})
	}
}

func Test_GetFullMessageByIDHandler(t *testing.T) {
	testMessage := model.FullMessage{ID: 1}

	tests := []struct {
		name            string
		mock            func(messageSrv *mocks.MessageService)
		input           string
		wantErr         bool
		expectedErr     lib.HttpError
		expectedMessage model.FullMessage
		expectedCode    int
	}{
		{
			name: "Ok: [message found]",
			mock: func(messageSrv *mocks.MessageService) {
				messageSrv.On("GetFullMessageByID", 1).Return(&testMessage, nil)
			},
			input:           "1",
			expectedMessage: testMessage,
			expectedCode:    http.StatusOK,
		},
		{
			name: "Error: [message not found]",
			mock: func(messageSrv *mocks.MessageService) {
				messageSrv.On("GetFullMessageByID", 1).Return(nil, fmt.Errorf("[Message] srv.GetFullMessageByID error: %w", pg.ErrFullMessageNotFound))
			},
			input:        "1",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 404, Name: "Not Found", Message: "full message not found"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "Error: [some internal error]",
			mock: func(messageSrv *mocks.MessageService) {
				messageSrv.On("GetFullMessageByID", 1).Return(nil, fmt.Errorf("[Message] srv.GetFullMessageByID error: some error"))
			},
			input:        "1",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 500, Name: "Internal Server Error", Message: "[Message] srv.GetFullMessageByID error: some error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "Error: [message id is not valid]",
			mock:         func(messageSrv *mocks.MessageService) {},
			input:        "Hello",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 400, Name: "Bad Request", Message: "message id is not valid"},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/message/%s", tt.input), nil)
			if err != nil {
				t.Fatalf("could not create request: %s", err)
			}

			rr := httptest.NewRecorder()

			log := logger.Get("debug")

			messageSrv := &mocks.MessageService{}
			tt.mock(messageSrv)

			handler := handler.New(&service.Manager{Message: messageSrv}, log)

			router := mux.NewRouter()
			router.HandleFunc("/message/{message_id}", handler.GetFullMessageByIDHandler)
			router.ServeHTTP(rr, req)

			decodedMessage := model.FullMessage{}
			decodedErr := lib.HttpError{}

			if tt.wantErr {
				json.NewDecoder(rr.Body).Decode(&decodedErr)

				assert.EqualValues(t, tt.expectedErr, decodedErr)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			} else {
				json.NewDecoder(rr.Body).Decode(&decodedMessage)

				assert.EqualValues(t, tt.expectedMessage, decodedMessage)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			}
		})
	}
}
