package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/VladPetriv/scanner_backend_api/internal/store/pg"
)

func (h *Handler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.log.Error("failed to get user id", zap.Error(err))

		h.WriteError(w, http.StatusBadRequest, "user id is not valid")

		return
	}

	user, err := h.service.User.GetUserByID(userID)
	if err != nil {
		h.log.Error("get user by id error", zap.String("id", strconv.Itoa(userID)), zap.Error(err))

		if errors.Is(err, pg.ErrUserNotFound) {
			err = errors.Unwrap(err)

			h.WriteError(w, http.StatusNotFound, err.Error())

			return
		}

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusOK, user)
}
