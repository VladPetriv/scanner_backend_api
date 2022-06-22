package handler

import (
	"fmt"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/VladPetriv/scanner_backend_api/internal/service"
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

	h.logAllRoutes(router)

	return router
}

func (h *Handler) logAllRoutes(router *mux.Router) {
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, err := route.GetPathTemplate()
		if err != nil {
			h.log.Error("", zap.Error(err))
		}

		met, err := route.GetMethods()
		if err != nil {
			h.log.Error("", zap.Error(err))
		}

		h.log.Info(fmt.Sprintf("Route - %s %s", tpl, met))

		return nil
	})
}
