package websocket

import (
	"encoding/json"
	"game/internal/manager/entity"
	"game/pkg/context"
	"net/http"
	"nhooyr.io/websocket"
)

type communication[E any] struct {
	conn *websocket.Conn
}

func NewCommunication[E any](conn *websocket.Conn) entity.Communication[E] {
	return &communication[E]{
		conn: conn,
	}
}

func (c communication[E]) Read(ctx context.Context) (E, error) {
	var e E
	_, b, err := c.conn.Read(ctx.NewContext())
	if err != nil {
		return e, err
	}

	err = json.Unmarshal(b, &e)
	if err != nil {
		return e, err
	}
	return e, nil
}

func (c communication[E]) Write(ctx context.Context, event E) error {
	b, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	return c.conn.Write(ctx.NewContext(), websocket.MessageBinary, b)
}

func (c communication[E]) Close(ctx context.Context) error {
	// TODO
	return c.conn.Close(http.StatusNoContent, "TODO")
}
