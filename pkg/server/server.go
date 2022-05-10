package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coolyrat/kit/pkg/config"
	"github.com/coolyrat/kit/pkg/server/app"
)

// TODO: logging substitute

type Server struct {
	*app.App
	http            *http.Server
	quit            chan os.Signal
	gracefulTimeout time.Duration
}

func New(app *app.App, handler http.Handler) *Server {
	return &Server{
		App: app,
		http: &http.Server{
			Addr:    config.GetString(config.PathServerPort),
			Handler: handler,
		},
		quit:            make(chan os.Signal),
		gracefulTimeout: config.GetDuration(config.PathGracefulTimeout),
	}
}

func (svr *Server) Run() {
	go func() {
		if err := svr.http.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Println("Server closed under request")
			} else {
				log.Fatal("Server closed unexpect: ", err)
			}
		}
	}()

	log.Println("Server started...")
	signal.Notify(svr.quit, syscall.SIGINT, syscall.SIGTERM)
	<-svr.quit

	svr.stop()
}

func (svr *Server) AppContext() context.Context {
	return svr.Context
}

func (svr *Server) stop() {
	log.Println("Shutting down server...")

	if err := svr.http.Shutdown(svr.Context); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server stopped...")
	svr.Stop()
	// TODO: graceful shutdown
	// time.Sleep(25 * time.Second)
}
