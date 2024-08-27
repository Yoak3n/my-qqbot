package hub

import (
	"my-qqbot/model"
	"time"
)

var (
	hub *model.EventHub
)

func runEventCircle() {
	for {
		for _, event := range hub.Pool {
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
