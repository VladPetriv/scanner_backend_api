package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/VladPetriv/scanner_backend_api/internal/store/pg"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (h *Handler) GetMessagesCountHandler(w http.ResponseWriter, r *http.Request) {
	count, err := h.service.Message.GetMessagesCount()
	if err != nil {
		h.log.Error("get messages count error", zap.Error(err))

		err = errors.Unwrap(err)

		if errors.Is(err, pg.ErrMessagesCountNotFound) {
			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, map[string]int{"count": count})
}

func (h *Handler) GetMessagesCountByChannelIDHandler(w http.ResponseWriter, r *http.Request) {
	channelID, err := strconv.Atoi(mux.Vars(r)["channel_id"])
	if err != nil {
		h.log.Error("get channel id from request error", zap.Error(err))

		h.WriteError(w, http.StatusBadRequest, "channel id is not valid")

		return
	}

	count, err := h.service.Message.GetMessagesCountByChannelID(channelID)
	if err != nil {
		h.log.Error("get messages count by channel id error", zap.String("id", strconv.Itoa(channelID)), zap.Error(err))

		err = errors.Unwrap(err)

		if errors.Is(err, pg.ErrMessagesCountNotFound) {
			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, map[string]int{"count": count})
}

func (h *Handler) GetFullMessagesByPageHandler(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		h.log.Error("get page from query error", zap.Error(err))

		h.WriteError(w, http.StatusBadRequest, "page is not valid")

		return
	}

	messages, err := h.service.Message.GetFullMessagesByPage(page)
	if err != nil {
		h.log.Error("get messages by page error", zap.String("page", strconv.Itoa(page)), zap.Error(err))

		err = errors.Unwrap(err)

		if errors.Is(err, pg.ErrFullMessagesNotFound) {
			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, messages)
}

func (h *Handler) GetFullMessagesByChannelIDAndPageHandler(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		h.log.Error("get page from query error", zap.Error(err))

		h.WriteError(w, http.StatusBadRequest, "page is not valid")

		return
	}

	channelID, err := strconv.Atoi(mux.Vars(r)["channel_id"])
	if err != nil {
		h.log.Error("get page from query error", zap.Error(err))

		h.WriteError(w, http.StatusBadRequest, "page is not valid")

		return
	}

	messages, err := h.service.Message.GetFullMessagesByChannelIDAndPage(channelID, page)
	if err != nil {
		h.log.Error("get full messages by page and channel id error", zap.String("id,page", fmt.Sprintf("%d,%d", channelID, page)))

		err = errors.Unwrap(err)

		if errors.Is(err, pg.ErrFullMessagesNotFound) {
			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, messages)
}

func (h *Handler) GetFullMessagesByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		h.log.Error("get user id from url error", zap.Error(err))

		h.WriteError(w, http.StatusBadRequest, "user id is not valid")

		return
	}

	messages, err := h.service.Message.GetFullMessagesByUserID(userID)
	if err != nil {
		h.log.Error("get full messages by user id error", zap.String("id", strconv.Itoa(userID)), zap.Error(err))

		err = errors.Unwrap(err)

		if errors.Is(err, pg.ErrFullMessagesNotFound) {
			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}
		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, messages)
}

func (h *Handler) GetFullMessageByIDHandler(w http.ResponseWriter, r *http.Request) {
	messageID, err := strconv.Atoi(mux.Vars(r)["message_id"])
	if err != nil {
		h.log.Error("get message id from url error", zap.Error(err))

		h.WriteError(w, http.StatusBadRequest, "message id is not valid")

		return
	}

	message, err := h.service.Message.GetFullMessageByID(messageID)
	if err != nil {
		h.log.Error("get message by id error", zap.String("id", strconv.Itoa(messageID)), zap.Error(err))

		err = errors.Unwrap(err)

		if errors.Is(err, pg.ErrFullMessageNotFound) {
			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, message)
}
