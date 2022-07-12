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

func Test_GetChannelsCountHandler(t *testing.T) {
	tests := []struct {
		name          string
		mock          func(channelSrv *mocks.ChannelService)
		wantErr       bool
		expectedErr   lib.HttpError
		expectedCount int
		expectedCode  int
	}{
		{
			name: "Ok: [channel count found]",
			mock: func(channelSrv *mocks.ChannelService) {
				channelSrv.On("GetChannelsCount").Return(10, nil)
			},
			expectedCount: 10,
			expectedCode:  http.StatusOK,
		},
		{
			name: "Ok: [channels count not found]",
			mock: func(channelSrv *mocks.ChannelService) {
				channelSrv.On("GetChannelsCount").Return(0, fmt.Errorf("[Channel] srv.GetChannelsCount error: %w", pg.ErrChannelsCountNotFound))
			},
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 404, Name: "Not Found", Message: "channels count not found"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "Error: [some internal error]",
			mock: func(channelSrv *mocks.ChannelService) {
				channelSrv.On("GetChannelsCount").Return(0, fmt.Errorf("[Channel] srv.GetChannelsCount error: some error"))
			},
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 500, Name: "Internal Server Error", Message: "[Channel] srv.GetChannelsCount error: some error"},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/channel/count", nil)
			if err != nil {
				t.Fatalf("could not create request: %s", err)
			}

			rr := httptest.NewRecorder()

			log := logger.Get("debug")

			channelSrv := &mocks.ChannelService{}
			tt.mock(channelSrv)

			handler := handler.New(&service.Manager{Channel: channelSrv}, log)

			router := mux.NewRouter()
			router.HandleFunc("/channel/count", handler.GetChannelsCountHandler)
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

			channelSrv.AssertExpectations(t)
		})
	}
}

func Test_GetChannelByNameHandler(t *testing.T) {
	testChannel := model.Channel{
		ID:       1,
		Name:     "test",
		Title:    "test test",
		ImageURL: "test.jpg",
	}

	tests := []struct {
		name            string
		mock            func(channelSrv *mocks.ChannelService)
		input           string
		wantErr         bool
		expectedErr     lib.HttpError
		expectedChannel model.Channel
		expectedCode    int
	}{
		{
			name: "Error: [channel found]",
			mock: func(channelSrv *mocks.ChannelService) {
				channelSrv.On("GetChannelByName", "test").Return(&testChannel, nil)
			},
			input:           "test",
			expectedChannel: testChannel,
			expectedCode:    http.StatusOK,
		},
		{
			name: "Ok: [channel not found]",
			mock: func(channelSrv *mocks.ChannelService) {
				channelSrv.On("GetChannelByName", "test").Return(nil, fmt.Errorf("[Channel] srv.GetChannelByName error: %w", pg.ErrChannelNotFound))
			},
			input:        "test",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 404, Name: "Not Found", Message: "channel not found"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "Error: [some internal error]",
			mock: func(channelSrv *mocks.ChannelService) {
				channelSrv.On("GetChannelByName", "test").Return(nil, fmt.Errorf("[Channel] srv.GetChannelByName error: some error"))
			},
			input:        "test",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 500, Name: "Internal Server Error", Message: "[Channel] srv.GetChannelByName error: some error"},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/channel/%s", tt.input), nil)
			if err != nil {
				t.Logf("could not create request: %s", err)
			}

			rr := httptest.NewRecorder()

			log := logger.Get("debug")

			channelSrv := &mocks.ChannelService{}
			tt.mock(channelSrv)

			handler := handler.New(&service.Manager{Channel: channelSrv}, log)

			router := mux.NewRouter()
			router.HandleFunc("/channel/{name}", handler.GetChannelByNameHandler)
			router.ServeHTTP(rr, req)

			decodedErr := lib.HttpError{}
			decodedChannel := model.Channel{}

			if tt.wantErr {
				json.NewDecoder(rr.Body).Decode(&decodedErr)

				assert.EqualValues(t, tt.expectedErr, decodedErr)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			} else {
				json.NewDecoder(rr.Body).Decode(&decodedChannel)

				assert.EqualValues(t, tt.expectedChannel, decodedChannel)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			}

			channelSrv.AssertExpectations(t)
		})
	}
}

func Test_GetChannelsByPageHandler(t *testing.T) {
	testChannels := []model.Channel{
		{ID: 1, Name: "test", Title: "test test", ImageURL: "test.jpg"},
		{ID: 2, Name: "test2", Title: "test2 test2", ImageURL: "test2.jpg"},
	}

	tests := []struct {
		name             string
		mock             func(channelSrv *mocks.ChannelService)
		input            string
		wantErr          bool
		expectedErr      lib.HttpError
		expectedChannels []model.Channel
		expectedCode     int
	}{
		{
			name: "Ok: [channels found]",
			mock: func(channelSrv *mocks.ChannelService) {
				channelSrv.On("GetChannelsByPage", 1).Return(testChannels, nil)
			},
			input:            "1",
			expectedChannels: testChannels,
			expectedCode:     http.StatusOK,
		},
		{
			name: "Error: [channels not found]",
			mock: func(channelSrv *mocks.ChannelService) {
				channelSrv.On("GetChannelsByPage", 1).Return(nil, fmt.Errorf("[Channel] srv.GetChannelsByPage error: %w", pg.ErrChannelsNotFound))
			},
			input:        "1",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 404, Name: "Not Found", Message: "channels not found"},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "Error: [some internal error]",
			mock: func(channelSrv *mocks.ChannelService) {
				channelSrv.On("GetChannelsByPage", 1).Return(nil, fmt.Errorf("[Channel] srv.GetChannelsByPage error: some error"))
			},
			input:        "1",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 500, Name: "Internal Server Error", Message: "[Channel] srv.GetChannelsByPage error: some error"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "Error: [page is not valid]",
			mock:         func(channelSrv *mocks.ChannelService) {},
			input:        "hello",
			wantErr:      true,
			expectedErr:  lib.HttpError{Code: 400, Name: "Bad Request", Message: "page is not valid"},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/channel/?page=%s", tt.input), nil)
			if err != nil {
				t.Fatalf("could not create request: %s", err)
			}

			rr := httptest.NewRecorder()

			log := logger.Get("debug")

			channelSrv := &mocks.ChannelService{}
			tt.mock(channelSrv)

			handler := handler.New(&service.Manager{Channel: channelSrv}, log)

			router := mux.NewRouter()
			router.HandleFunc("/channel/", handler.GetChannelsByPageHandler)
			router.ServeHTTP(rr, req)

			decodedChannels := []model.Channel{}
			decodedErr := lib.HttpError{}

			if tt.wantErr {
				json.NewDecoder(rr.Body).Decode(&decodedErr)

				assert.EqualValues(t, tt.expectedErr, decodedErr)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			} else {
				json.NewDecoder(rr.Body).Decode(&decodedChannels)

				assert.EqualValues(t, tt.expectedChannels, decodedChannels)
				assert.EqualValues(t, tt.expectedCode, rr.Code)
			}

			channelSrv.AssertExpectations(t)
		})
	}
}
