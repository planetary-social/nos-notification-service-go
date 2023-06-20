package domain

import (
	"bytes"

	"github.com/boreq/errors"
	"github.com/nbd-wtf/go-nostr"
)

type Event struct {
	pubKey  PublicKey
	content []byte
}

func NewEventFromEnvelope(envelope nostr.EventEnvelope) (Event, error) {
	ok, err := envelope.CheckSignature()
	if err != nil {
		return Event{}, errors.Wrap(err, "error checking signature")
	}

	if !ok {
		return Event{}, errors.New("invalid signature")
	}

	pubKey, err := NewPublicKey(envelope.PubKey)
	if err != nil {
		return Event{}, errors.Wrap(err, "error creating a pub key")
	}

	return Event{
		pubKey:  pubKey,
		content: []byte(envelope.Content),
	}, nil
}

func (e Event) PubKey() PublicKey {
	return e.pubKey
}

func (e Event) Content() []byte {
	return bytes.Clone(e.content)
}
