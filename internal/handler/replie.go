package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/VladPetriv/scanner_backend_api/internal/store/pg"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (h *Handler) GetFullRepliesByMessageIDHandler(w http.ResponseWriter, r *http.Request) {
	messageID, err := strconv.Atoi(mux.Vars(r)["message_id"])
	if err != nil {
		h.log.Error("failed to get message id", zap.Error(err))

		h.WriteError(w, http.StatusBadRequest, "message id is not valid")

		return
	}

	replies, err := h.service.Replie.GetFullRepliesByMessageID(messageID)
	if err != nil {
		h.log.Error("get full replies by message id error", zap.String("id", strconv.Itoa(messageID)), zap.Error(err))

		err = errors.Unwrap(err)

		if errors.Is(err, pg.ErrFullRepliesNotFound) {
			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, replies)
}
