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

func (k EventKind) Int() int {
	return k.k
}
