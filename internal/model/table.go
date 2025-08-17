package model

type FromTable struct {
	Id      int64
	Private bool
}
type LiveTable struct {
	Id   int64
	From FromTable
}
