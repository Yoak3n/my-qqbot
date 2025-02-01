package hub

import (
	"my-qqbot/internal/model"
	"time"
)

type eventHub struct {
	Pool  []*model.Event
	Begin bool
}

var (
	hub *eventHub
)

func (h *eventHub) runEventCircle() {
	for {
		for _, event := range h.Pool {
			if event.Timer != nil && !event.Running {
				go func() {
					if event.Running {
						return
					}
					event.Running = true
					<-event.Timer.C
					event.Action(event.From)
					event.Timer.Reset(24 * time.Hour)
					event.Running = false
				}()
			}
		}
	}
}
