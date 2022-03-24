package rest

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type builder struct {
	port            int
	handler         http.Handler
	cleanupFunc     func()
	gracefulTimeout time.Duration
}

// NewServerBuilder defines a new server builder.
func NewServerBuilder() *builder {
	return &builder{
		port:            3000,
		gracefulTimeout: time.Second * 30,
	}
}

func (b *builder) WithPort(p int) *builder {
	b.port = p
	return b
}

func (b *builder) WithCleanup(fn func()) *builder {
	b.cleanupFunc = fn
	return b
}

func (b *builder) WithHandler(h http.Handler) *builder {
	b.handler = h
	return b
}

func (b *builder) WithShutdownTimeout(d time.Duration) *builder {
	b.gracefulTimeout = d
	return b
}

func (b *builder) Build() *server {
	return &server{
		http: &http.Server{
			Addr:    fmt.Sprintf(":%d", b.port),
			Handler: b.handler,
		},
		cleanup:         b.cleanupFunc,
		quit:            make(chan os.Signal, 1),
		gracefulTimeout: b.gracefulTimeout,
	}
}
