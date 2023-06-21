package domain

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/boreq/errors"
)

type EventId struct {
	s string
}

func NewEventId(s string) (EventId, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return EventId{}, errors.Wrap(err, "error decoding hex")
	}

	if len(b) != sha256.Size {
		return EventId{}, errors.New("invalid length")
	}

	s = hex.EncodeToString(b)
	return EventId{s}, nil
}

func (id EventId) Hex() string {
	return id.s
}

func (id EventId) Bytes() []byte {
	b, err := hex.DecodeString(id.s)
	if err != nil {
		panic(err)
	}
	return b
}
