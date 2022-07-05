package handler

import (
	"net/http"
	"strings"
)

func (h *Handler) AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			h.WriteError(w, http.StatusUnauthorized, "user is not authorized")

			return
		}
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			h.WriteError(w, http.StatusUnauthorized, "invalid auth header")
			return
		}

		if len(headerParts[1]) == 0 {
			h.WriteError(w, http.StatusUnauthorized, "token is empty")
			return
		}

		userEmail, err := h.service.Jwt.ParseToken(headerParts[1])
		if err != nil {
			h.WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}

		r.Header.Set("email", userEmail)

		next.ServeHTTP(w, r)
	})
}
