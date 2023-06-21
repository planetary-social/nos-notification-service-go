package domain

import "github.com/boreq/errors"

type EventTag struct {
	tag []string
}

func NewEventTag(tag []string) (EventTag, error) {
	if len(tag) < 2 {
		return EventTag{}, errors.New("tag needs at least two fields I recon")
	}

	return EventTag{tag}, nil
}
