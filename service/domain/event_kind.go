package domain

import "github.com/boreq/errors"

type EventKind struct {
	k int
}

func NewEventKind(k int) (EventKind, error) {
	if k < 0 {
		return EventKind{}, errors.New("kind must be positive")
	}
	return EventKind{k}, nil
}

func MustNewEventKind(k int) EventKind {
	v, err := NewEventKind(k)
	if err != nil {
		panic(err)
	}
	return v
}

func (k EventKind) Int() int {
	return k.k
}
