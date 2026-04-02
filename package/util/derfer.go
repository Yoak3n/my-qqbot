package util

import "my-qqbot/internal/hub"

func DeferCall(preblock func(), f func()) {
	hub.SendMessage()
}
