package kitsqs

import (
	"context"
)

type Context struct {
	context.Context
	abortErr error
}

func (c *Context) Abort(err error) {
	c.abortErr = err
}

func (c *Context) Next() {
	// c.index++
	// for c.index < int8(len(c.handlers)) {
	// 	c.handlers[c.index](c)
	// 	c.index++
	// }
}
