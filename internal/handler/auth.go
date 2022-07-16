package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/VladPetriv/scanner_backend_api/internal/model"
	"github.com/VladPetriv/scanner_backend_api/internal/store/pg"
	"github.com/VladPetriv/scanner_backend_api/pkg/utils"
)

// SignUpHandler godoc
// @ID           create-user
// @Summary      sign-up
// @Description  Handler will create new user and return message
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      model.WebUser  true  "user info"
// @Success      200    {string}  string         "user created"
// @Failure      400    {object}  lib.HttpError  "bad request"
// @Failure      409    {object}  lib.HttpError  "user with email is exist"
// @Failure      500    {object}  lib.HttpError  "internal server error"
// @Router       /auth/sign-up [post]
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

	err = h.service.WebUser.CreateWebUser(&user)
	if err != nil {
		h.log.Error("failed to create user", zap.Any("user structure", user), zap.Error(err))

		h.WriteError(w, http.StatusInternalServerError, err.Error())

		return
	}

	h.WriteJSON(w, http.StatusCreated, "user created")
}

// SignInHandler godoc
// @ID           login-user
// @Summary      sigi-up
// @Description  Handler will login user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      model.WebUser  true  "user info"
// @Success      200    {string}  string         "token"
// @Failure      400    {object}  lib.HttpError  "bad request"
// @Failure      404    {object}  lib.HttpError  "user with this email not found"
// @Failure      500    {object}  lib.HttpError  "internal server error"
// @Router       /auth/sign-in [post]
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

			h.WriteError(w, http.StatusNotFound, "user not found")

			return
		}
	}

	ok := utils.ComparePassword(user.Password, candidate.Password)
	if !ok {
		h.WriteError(w, http.StatusUnauthorized, "password is incorrect")

		return
	}

	token, err := h.service.Jwt.GenerateToken(candidate.Email)
	if err != nil {
		h.log.Error("failed to generate error", zap.Error(err))

		h.WriteError(w, http.StatusInternalServerError, "failed to generate jwt token")

		return
	}

	h.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}
