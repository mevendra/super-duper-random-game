package context

import (
	"context"
	"fmt"
)

type Context interface {
	NewContext() context.Context
	Debug(format string, args ...any)
}

type ctx struct {
	context.Context
}

func New() Context {
	return &ctx{context.Background()}
}

func (c ctx) NewContext() context.Context {
	return context.WithoutCancel(c.Context)
}

func (c ctx) Debug(format string, args ...any) {
	fmt.Print("DEBUG: ")
	fmt.Printf(format, args...)
	fmt.Println()
}
