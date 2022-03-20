package rest

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type server struct {
	cleanup         func()
	http            *http.Server
	quit            chan os.Signal
	gracefulTimeout time.Duration
}

func (svr *server) Start() {

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

	svr.Stop()
}

func (svr *server) Stop() {
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), svr.gracefulTimeout)
	defer cancel()
	if err := svr.http.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	if svr.cleanup != nil {
		log.Println("Server resources cleanup...")
		svr.cleanup()
	}

	log.Println("Server exiting...")
}
