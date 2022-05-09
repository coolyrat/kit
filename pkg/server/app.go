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
	"github.com/gin-gonic/gin"
)

type App struct {
	context.Context
	http            *http.Server
	quit            chan os.Signal
	gracefulTimeout time.Duration
}

func NewApp(ctx context.Context, e *gin.Engine) *App {
	return &App{
		Context: ctx,
		http: &http.Server{
			Addr:    config.GetString(config.PathServerPort),
			Handler: defaultHandler(e),
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

func (app *App) stop() {
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(app.Context, app.gracefulTimeout)
	// defer cancel()
	if err := app.http.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server stopped...")
	cancel()
	time.Sleep(25 * time.Second)
}
