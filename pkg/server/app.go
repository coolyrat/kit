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
	http            *http.Server
	quit            chan os.Signal
	gracefulTimeout time.Duration
}

func NewApp(e *gin.Engine) *App {
	return &App{

		http: &http.Server{
			Addr:    config.GetString(config.PathServerPort),
			Handler: defaultHandler(e),
		},
		quit:            make(chan os.Signal),
		gracefulTimeout: time.Second * 20,
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

	ctx, cancel := context.WithTimeout(context.Background(), app.gracefulTimeout)
	defer cancel()
	if err := app.http.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server stopped...")
}
