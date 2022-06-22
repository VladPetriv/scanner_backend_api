package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/VladPetriv/scanner_backend_api/pkg/config"
)

var (
	readTimeout  = time.Second * 10
	writeTimeout = time.Second * 10
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Start(handler http.Handler) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	return s.httpServer.ListenAndServe()
}
