package types

import (
	"fmt"
	"time"
)

type ServerOptions struct {
	Name        string
	Host        string
	GrpcPort    int
	HttpPort    int
	UseSSL      bool
	ReflectGRPC bool

	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
}

func (o ServerOptions) GrpcAddr() string {
	return fmt.Sprintf("%s:%d", o.Host, o.GrpcPort)
}

func (o ServerOptions) HttpAddr() string {
	return fmt.Sprintf("%s:%d", o.Host, o.HttpPort)
}

type ServerOptionsFunc = func(*ServerOptions)

func DefaultServerOptions() ServerOptions {
	return ServerOptions{
		Name:        "server",
		Host:        "",
		GrpcPort:    3000,
		HttpPort:    4000,
		UseSSL:      false,
		ReflectGRPC: true,

		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}
}

func NewServerOptions(fns ...ServerOptionsFunc) ServerOptions {
	options := DefaultServerOptions()
	for _, fn := range fns {
		fn(&options)
	}

	return options
}
