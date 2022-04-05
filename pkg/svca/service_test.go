package svca

import (
	"testing"
	"time"

	"github.com/coolyrat/kit/pkg/config"
	"github.com/coolyrat/kit/pkg/logr"
)

func TestService(t *testing.T) {
	NewService()
	for {
		config.Print()
		logr.Info("hello")
		time.Sleep(15 * time.Second)
	}
}
