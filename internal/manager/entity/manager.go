package entity

import (
	"game/pkg/context"
	"sync"
)

type manager[E any, C Communication[E]] struct {
}

func New[E any, C Communication[E]]() Manager[E, C[E]] {
	return &manager[E, C]{}
}

func (m *manager[E, C]) Start(size int, ch <-chan C, commManager func() C) {
	ctx := context.New()
	ctx.Debug("Starting new manager")

	commList := make([]Communication[E], 0, size)
	for {
		comm := <-ch
		commList = append(commList, comm)
		ctx.Debug("New communication received")
		if l := len(commList); l < size {
			ctx.Debug("Missing %d more communications to start", size-l)
			continue
		}

		ctx.Debug("Starting new group")
		ncList := commList[:size]
		ncList = append(ncList, commManager())
		go m.start(ncList)
		commList = commList[size:]
	}
}

func (m *manager[E, C]) start(communications []Communication[E]) {
	ctx := context.New()
	wg := sync.WaitGroup{}
	ctx.Debug("Starting communication reader")
	for rI, rComm := range communications {
		wg.Add(1)
		comm := rComm
		i := rI
		go func() {
			e, err := comm.Read(ctx)
			if err != nil {
				ctx.Debug("Error reading from communication: %s", err.Error())
				err = comm.Close(ctx)
				if err != nil {
					ctx.Debug("Error closing communication: %s", err.Error())
				}
				wg.Done()
				return
			}

			ctx.Debug("Event received")
			m.notifyOthers(ctx, e, i, communications)
		}()
	}

	ctx.Debug("All communications have started")
	wg.Wait()
	ctx.Debug("All communications have stopped")
}

func (m *manager[E, C]) notifyOthers(ctx context.Context, event E, ignore int, communications []Communication[E]) {
	for i, rComm := range communications {
		if i == ignore {
			continue
		}

		comm := rComm
		go func() {
			err := comm.Write(ctx, event)
			if err != nil {
				ctx.Debug("Error notifying communication: %s", err.Error())
			}
		}()
	}
}
