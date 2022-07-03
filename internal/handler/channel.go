package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/VladPetriv/scanner_backend_api/internal/store/pg"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (h *Handler) GetChannelsCountHandler(w http.ResponseWriter, r *http.Request) {
	count, err := h.service.Channel.GetChannelsCount()
	if err != nil {
		h.log.Error("get channels count error", zap.Error(err))

		err = errors.Unwrap(err)

		if errors.Is(err, pg.ErrChannelsCountNotFound) {
			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, map[string]int{"count": count})
}

func (h *Handler) GetChannelByNameHandler(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	channel, err := h.service.Channel.GetChannelByName(name)
	if err != nil {
		h.log.Error("get channel by name error", zap.String("name", name), zap.Error(err))

		err = errors.Unwrap(err)

		if errors.Is(err, pg.ErrChannelNotFound) {
			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, channel)
}

func (h *Handler) GetChannelsByPageHandler(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		h.log.Error("failed to get page", zap.Error(err))

		h.WriteError(w, http.StatusBadRequest, "page is not valid")

		return
	}

	channels, err := h.service.Channel.GetChannelsByPage(page)
	if err != nil {
		h.log.Error("get channels by page error", zap.String("page", strconv.Itoa(page)), zap.Error(err))

		err = errors.Unwrap(err)

		if errors.Is(err, pg.ErrChannelsNotFound) {
			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, channels)
}
