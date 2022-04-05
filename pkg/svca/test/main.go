package main

import (
	"time"

	"github.com/coolyrat/kit/pkg/config"
	"github.com/coolyrat/kit/pkg/logr"
	"github.com/coolyrat/kit/pkg/svca"
)

func main() {
	svca.NewService()
	for {
		config.Print()
		logr.Info("hello")
		time.Sleep(15 * time.Second)
	}
}
