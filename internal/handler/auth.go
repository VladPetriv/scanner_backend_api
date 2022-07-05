package handler

import (
	"encoding/json"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"net/http"

	"github.com/VladPetriv/scanner_backend_api/internal/model"
	"github.com/VladPetriv/scanner_backend_api/internal/store/pg"
	"github.com/VladPetriv/scanner_backend_api/pkg/utils"
)

func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	user := model.WebUser{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.log.Error("failed to decode request body", zap.Error(err))

		h.WriteError(w, http.StatusBadRequest, "can't decode request body")

		return
	}

	candidate, err := h.service.WebUser.GetWebUserByEmail(user.Email)
	if err != nil {
		if !errors.Is(err, pg.ErrWebUserNotFound) {
			h.log.Error("failed to get candidate", zap.String("email", user.Email), zap.Error(err))

			h.WriteError(w, http.StatusInternalServerError, err.Error())

			return
		}
	}

	if candidate != nil {
		h.WriteError(w, http.StatusConflict, fmt.Sprintf("user with email %s is exist", user.Email))

		return
	}

	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		h.log.Error("failed to hash password", zap.Error(err))

		return
	}

	err = h.service.WebUser.CreateWebUser(&user)
	if err != nil {
		h.log.Error("failed to create user", zap.Any("user structure", user), zap.Error(err))

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusCreated, "user created")
}

func (h *Handler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	user := model.WebUser{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.log.Error("failed to decode request body", zap.Error(err))

		h.WriteError(w, http.StatusBadRequest, "can't decode request body")

		return
	}

	candidate, err := h.service.WebUser.GetWebUserByEmail(user.Email)
	if err != nil {
		if !errors.Is(err, pg.ErrWebUserNotFound) {
			h.log.Error("failed to get candidate", zap.String("email", user.Email), zap.Error(err))

			h.WriteError(w, http.StatusInternalServerError, err.Error())

			return
		}

		if errors.Is(err, pg.ErrWebUserNotFound) {
			h.log.Error("failed to get user", zap.String("email", user.Email), zap.Error(err))

			h.WriteError(w, http.StatusInternalServerError, "user not found")

			return
		}
	}

	token, err := h.service.Jwt.GenerateToken(candidate.Email)
	if err != nil {
		h.log.Error("failed to generate error", zap.Error(err))

		h.WriteError(w, http.StatusInternalServerError, "failed to generate jwt token")

		return
	}

	h.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}
