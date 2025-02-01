package model

type (
	Notification struct {
		Private bool
		Message string
		Target  int64
		Picture []string
	}
	From struct {
		Private bool
		Id      int64
	}
)
