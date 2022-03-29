package svca

import (
	"fmt"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	s := NewService()
	for {
		fmt.Println(s.config.Age)
		time.Sleep(5 * time.Second)
	}
}
