package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/VladPetriv/scanner_backend_api/internal/store/pg"
)

// GetFullRepliesByMessageIDHandler godoc
// @ID           get-full-replies-by-message-id
// @Summary      GetFullRepliesByMessageID
// @Description  Handler will return full replies by message id from url
// @Tags         replie
// @Produce      json
// @Param        message_id  path      integer           true  "message id"
// @Success      200         {array}   model.FullReplie  "replies by message id"
// @Failure      400         {object}  lib.HttpError     "bad request"
// @Failure      404         {object}  lib.HttpError     "replies not found"
// @Failure      500         {object}  lib.HttpError     "internal server error"
// @Router       /replie/{message_id} [get]
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

		if errors.Is(err, pg.ErrFullRepliesNotFound) {
			err = errors.Unwrap(err)

			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, replies)
}
