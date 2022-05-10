package app

import (
	"context"

	"github.com/coolyrat/kit/pkg/config"
)

// Version TODO: inject from build
var Version = "local"

type App struct {
	context.Context
	Stop context.CancelFunc

	Name    string
	Version string
	Mode    string
	Env     string
}

func New() *App {
	var app App
	if err := config.Unmarshal(config.PathApplication, &app); err != nil {
		// TODO
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	app.Context = ctx
	app.Stop = cancel
	app.Mode = Version
	return &app
}
