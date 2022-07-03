package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/VladPetriv/scanner_backend_api/internal/service"
	"github.com/VladPetriv/scanner_backend_api/pkg/lib"
	"github.com/VladPetriv/scanner_backend_api/pkg/logger"
)

type Handler struct {
	service *service.Manager
	log     *logger.Logger
}

func New(service *service.Manager, log *logger.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

func (h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()

	channel := router.PathPrefix("/channel").Subrouter()
	channel.HandleFunc("/count", h.GetChannelsCountHandler).Methods(http.MethodGet)
	channel.HandleFunc("/{name}", h.GetChannelByNameHandler).Methods(http.MethodGet)
	channel.HandleFunc("/", h.GetChannelsByPageHandler).Methods(http.MethodGet)

	user := router.PathPrefix("/user").Subrouter()
	user.HandleFunc("/{id}", h.GetUserByIDHandler).Methods(http.MethodGet)

	replie := router.PathPrefix("/replie").Subrouter()
	replie.HandleFunc("/{message_id}", h.GetFullRepliesByMessageIDHandler).Methods(http.MethodGet)

	h.logAllRoutes(router)

	return router
}

func (h *Handler) logAllRoutes(router *mux.Router) {
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, err := route.GetPathTemplate()
		if err != nil {
			h.log.Error("", zap.Error(err))
		}

		met, _ := route.GetMethods()

		h.log.Info(fmt.Sprintf("Route - %s %s", tpl, met))

		return nil
	})
}

func (h *Handler) WriteJSON(w http.ResponseWriter, httpCode int, data interface{}) {
	w.WriteHeader(httpCode)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (h *Handler) WriteError(w http.ResponseWriter, httpCode int, err string) {
	w.WriteHeader(httpCode)

	if err != "" {
		json.NewEncoder(w).Encode(lib.HttpError{
			Code:    httpCode,
			Name:    http.StatusText(httpCode),
			Message: err,
		})
	}
}
