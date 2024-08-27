package model

import "time"

type (
	Event struct {
		Name    string
		Action  func(from *From)
		From    *From
		Timer   *time.Timer
		Running bool
	}

	EventHub struct {
		Pool  []*Event
		Begin bool
	}
)
