package app

import "context"

type App struct {
	context.Context
	Stop context.CancelFunc
}

func New() *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{Context: ctx, Stop: cancel}
}
