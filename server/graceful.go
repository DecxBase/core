package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/phuslu/log"
	"google.golang.org/grpc"
)

type GracefulOptions struct {
	Logger  log.Logger
	Setup   func()
	Cleanup func()
	Closer  func()
	Start   func() error
	Stop    func(context.Context) error
}

func GracefulGrpc(srv *grpc.Server, addr string, l log.Logger, closers ...func()) (<-chan error, error) {
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err)
	}

	return GracefulRun(GracefulOptions{
		Logger: l,
		Start: func() error {
			l.Info().Msgf("Listening GRPC [%s]", addr)

			return srv.Serve(lis)
		},
		Stop: func(ctx context.Context) error {
			srv.GracefulStop()

			return nil
		},
		Closer: func() {
			err := lis.Close()
			if err != nil && !strings.Contains(err.Error(), "closed network") {
				l.Debug().Msgf("error closing gRPC listener: %v", err)
			}
		},
	}, closers...)
}

func GracefulHttp(srv *http.Server, l log.Logger, closers ...func()) (<-chan error, error) {
	return GracefulRun(GracefulOptions{
		Logger: l,
		Cleanup: func() {
			srv.SetKeepAlivesEnabled(false)
		},
		Start: func() error {
			l.Info().Msgf("Listening HTTP [%s]", srv.Addr)

			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				return err
			}

			return nil
		},
		Stop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	}, closers...)
}

func GracefulRun(opts GracefulOptions, closers ...func()) (<-chan error, error) {
	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGHUP,
	)

	go func() {
		<-ctx.Done()

		// opts.Logger.Info().Msg("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {
			for _, fn := range closers {
				fn()
			}

			if opts.Closer != nil {
				opts.Closer()
			}

			stop()
			cancel()
			close(errC)
		}()

		if opts.Cleanup != nil {
			opts.Cleanup()
		}

		if err := opts.Stop(ctxTimeout); err != nil {
			errC <- err
		}

		opts.Logger.Info().Msg("Shutdown completed")
	}()

	go func() {
		if err := opts.Start(); err != nil {
			errC <- err
		}
	}()

	return errC, nil
}
