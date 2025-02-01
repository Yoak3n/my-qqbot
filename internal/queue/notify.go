package queue

import (
	"my-qqbot/internal/model"
)

var Notify chan *model.Notification

func init() {
	Notify = make(chan *model.Notification, 100)
}
