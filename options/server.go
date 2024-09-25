package options

import (
	"time"

	"github.com/DecxBase/core/types"
)

func WithName(name string) types.ServerOptionsFunc {
	return func(o *types.ServerOptions) {
		o.Name = name
	}
}

func WithHost(host string) types.ServerOptionsFunc {
	return func(o *types.ServerOptions) {
		o.Host = host
	}
}

func WithGrpcPort(port int) types.ServerOptionsFunc {
	return func(o *types.ServerOptions) {
		o.GrpcPort = port
	}
}

func WithHttpPort(port int) types.ServerOptionsFunc {
	return func(o *types.ServerOptions) {
		o.HttpPort = port
	}
}

func WithSSL(useSSL bool) types.ServerOptionsFunc {
	return func(o *types.ServerOptions) {
		o.UseSSL = useSSL
	}
}

func WithReflectGRPC(r bool) types.ServerOptionsFunc {
	return func(o *types.ServerOptions) {
		o.ReflectGRPC = r
	}
}

func WithReadTimeout(t time.Duration) types.ServerOptionsFunc {
	return func(o *types.ServerOptions) {
		o.ReadTimeout = t
	}
}

func WithWriteTimeout(t time.Duration) types.ServerOptionsFunc {
	return func(o *types.ServerOptions) {
		o.WriteTimeout = t
	}
}

func WithIdleTimeout(t time.Duration) types.ServerOptionsFunc {
	return func(o *types.ServerOptions) {
		o.IdleTimeout = t
	}
}

func WithHandlerTimeout(t time.Duration) types.ServerOptionsFunc {
	return func(o *types.ServerOptions) {
		o.HandlerTimeout = t
	}
}

func WithReadHeaderTimeout(t time.Duration) types.ServerOptionsFunc {
	return func(o *types.ServerOptions) {
		o.ReadHeaderTimeout = t
	}
}

func WithEnvHost(key string) types.ServerOptionsFunc {
	return func(o *types.ServerOptions) {
		o.Host = ReadEnv(o.Name, key, o.Host)
	}
}

func WithEnvGrpcPort(key string) types.ServerOptionsFunc {
	return func(o *types.ServerOptions) {
		o.GrpcPort = ReadEnv(o.Name, key, o.GrpcPort)
	}
}

func WithEnvHttpPort(key string) types.ServerOptionsFunc {
	return func(o *types.ServerOptions) {
		o.HttpPort = ReadEnv(o.Name, key, o.HttpPort)
	}
}
