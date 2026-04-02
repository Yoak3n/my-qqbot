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
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		for _, event := range h.Pool {
			if event == nil || event.Timer == nil || event.Running {
				continue
			}
			event.Running = true
			go func(e *model.Event) {
				<-e.Timer.C
				e.Action(e.From)
				e.Timer.Reset(24 * time.Hour)
				e.Running = false
			}(event)
		}
	}
}
