package entity

import "game/pkg/context"

type Manager[E any, C Communication[E]] interface {
	Start(size int, ch <-chan C, commManager func() C)
}

type Communication[E any] interface {
	Read(ctx context.Context) (E, error)
	Write(ctx context.Context, event E) error
	Close(ctx context.Context) error
}
