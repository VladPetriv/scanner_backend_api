package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/VladPetriv/scanner_backend_api/internal/store/pg"
)

// GetUserByIDHandler godoc
// @ID           get-user-by-id
// @Summary      GetUserByID
// @Description  Handler will return user by id from url
// @Tags         user
// @Produce      json
// @Param        id   path      integer        true  "user id"
// @Success      200  {object}  model.User     "user by id"
// @Failure      400  {object}  lib.HttpError  "bad request"
// @Failure      404  {object}  lib.HttpError  "user not found"
// @Failure      500  {object}  lib.HttpError  "internal server error"
// @Router       /user/{id} [get]
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
