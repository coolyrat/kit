package server

import "context"

type Server interface {
	Run()
	AppContext() context.Context
}
