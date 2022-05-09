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

type App struct {
	*app.App
	http            *http.Server
	quit            chan os.Signal
	gracefulTimeout time.Duration
}

func NewApp(app *app.App, handler http.Handler) *App {
	return &App{
		App: app,
		http: &http.Server{
			Addr:    config.GetString(config.PathServerPort),
			Handler: handler,
		},
		quit:            make(chan os.Signal),
		gracefulTimeout: config.GetDuration(config.PathGracefulTimeout),
	}
}

func (app *App) Run() {
	go func() {
		if err := app.http.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Println("Server closed under request")
			} else {
				log.Fatal("Server closed unexpect: ", err)
			}
		}
	}()

	log.Println("Server started...")
	signal.Notify(app.quit, syscall.SIGINT, syscall.SIGTERM)
	<-app.quit

	app.stop()
}

func (app *App) AppContext() context.Context {
	return app.Context
}

func (app *App) stop() {
	log.Println("Shutting down server...")

	// defer cancel()
	if err := app.http.Shutdown(app.Context); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server stopped...")
	app.Stop()
	time.Sleep(25 * time.Second)
}
