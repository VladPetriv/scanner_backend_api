package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/VladPetriv/scanner_backend_api/internal/model"
	"github.com/VladPetriv/scanner_backend_api/internal/store/pg"
)

func (h *Handler) GetSavedMessagesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		h.log.Error("failed to get user id from Request", zap.Error(err))

		h.WriteError(w, http.StatusBadRequest, "user id is not valid")

		return
	}

	messages, err := h.service.Saved.GetSavedMessages(userID)
	if err != nil {
		h.log.Error("failed to get saved messages", zap.String("user id", strconv.Itoa(userID)), zap.Error(err))

		if errors.Is(err, pg.ErrSavedMessagesNotFound) {
			err = errors.Unwrap(err)

			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, messages)
}

func (h *Handler) CreateSavedMessageHandler(w http.ResponseWriter, r *http.Request) {
	saved := model.Saved{}
	if err := json.NewDecoder(r.Body).Decode(&saved); err != nil {
		h.log.Error("failed to decode request body", zap.Error(err))

		h.WriteError(w, http.StatusBadRequest, "can't decode request body")

		return
	}

	err := h.service.Saved.CreateSavedMessage(&saved)
	if err != nil {
		h.log.Error("failed to create saved message", zap.Any("saved structure", saved), zap.Error(err))

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusCreated, "saved message created")
}

func (h *Handler) DeleteSavedMessageHandler(w http.ResponseWriter, r *http.Request) {
	messageID, err := strconv.Atoi(mux.Vars(r)["message_id"])
	if err != nil {
		h.log.Error("failed to get message id", zap.Error(err))

		h.WriteError(w, http.StatusBadRequest, "message id is not valid")

		return
	}

	err = h.service.Saved.DeleteSavedMessage(messageID)
	if err != nil {
		h.log.Error("failed to delete saved message", zap.String("message id", strconv.Itoa(messageID)), zap.Error(err))

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, "saved message deleted")
}
