package models

type Jump struct {
	Start int
	End   int
}

func (j Jump) IsSnake() bool {
	return j.End < j.Start
}
