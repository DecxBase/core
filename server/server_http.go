package server

import (
	"errors"
	"net/http"

	"github.com/DecxBase/core/logger"
	"github.com/DecxBase/core/utils"
	"github.com/phuslu/log"
)

func (s composedServer) BuildHttp(l log.Logger) (*http.ServeMux, error) {
	router := http.NewServeMux()

	for _, hnd := range s.httpHandlers {
		hnd.SetLogger(logger.CreateFrom(l, "service", hnd.HttpIdentifier()))
		hnd.RegisterRoutes(router)
	}

	router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteError(w, http.StatusNotFound, errors.New("api endpoint not found"))
	})

	return router, nil
}

func (s composedServer) RunHttp(closers ...func()) error {
	if len(s.httpHandlers) < 1 {
		return errors.New("no http handlers registered")
	}

	l := logger.CreateFrom(s.Logger, "type", "http")

	router, err := s.BuildHttp(l)
	if err != nil {
		return err
	}

	srv := &http.Server{
		ReadTimeout:       s.opts.ReadTimeout,
		WriteTimeout:      s.opts.WriteTimeout,
		IdleTimeout:       s.opts.IdleTimeout,
		ReadHeaderTimeout: s.opts.ReadHeaderTimeout,

		Addr:    s.opts.HttpAddr(),
		Handler: router,
		// Handler: http.TimeoutHandler(router, s.opts.HandlerTimeout, "Timeout\n"),
	}

	errC, err := GracefulHttp(srv, l)
	if err != nil {
		return err
	}

	if err := <-errC; err != nil {
		return err
	}

	return nil
}
