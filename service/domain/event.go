package domain

import (
	"time"

	"github.com/boreq/errors"
	"github.com/nbd-wtf/go-nostr"
)

type Event struct {
	id        EventId
	pubKey    PublicKey
	createdAt time.Time
	kind      EventKind
	tags      []EventTag
	content   string
	sig       EventSignature

	libevent nostr.Event
}

func NewEvent(libevent nostr.Event) (Event, error) {
	ok, err := libevent.CheckSignature()
	if err != nil {
		return Event{}, errors.Wrap(err, "error checking signature")
	}

	if !ok {
		return Event{}, errors.New("invalid signature")
	}

	id, err := NewEventId(libevent.ID)
	if err != nil {
		return Event{}, errors.Wrap(err, "error creating an event id")
	}

	pubKey, err := NewPublicKeyFromHex(libevent.PubKey)
	if err != nil {
		return Event{}, errors.Wrap(err, "error creating a pub key")
	}

	createdAt := time.Unix(int64(libevent.CreatedAt), 0).UTC()

	kind, err := NewEventKind(libevent.Kind)
	if err != nil {
		return Event{}, errors.Wrap(err, "error creating event kind")
	}

	var tags []EventTag
	for _, libtag := range libevent.Tags {
		eventTag, err := NewEventTag(libtag)
		if err != nil {
			return Event{}, errors.Wrap(err, "error creating a tag")
		}
		tags = append(tags, eventTag)
	}

	sig, err := NewEventSignature(libevent.Sig)
	if err != nil {
		return Event{}, errors.Wrap(err, "error creating a signature")
	}

	return Event{
		id:        id,
		pubKey:    pubKey,
		createdAt: createdAt,
		kind:      kind,
		tags:      tags,
		content:   libevent.Content,
		sig:       sig,

		libevent: libevent,
	}, nil
}

func (e Event) Id() EventId {
	return e.id
}

func (e Event) PubKey() PublicKey {
	return e.pubKey
}

func (e Event) CreatedAt() time.Time {
	return e.createdAt
}

func (e Event) Kind() EventKind {
	return e.kind
}

func (e Event) Tags() []EventTag {
	return e.tags
}

func (e Event) Content() string {
	return e.content
}

func (e Event) Sig() EventSignature {
	return e.sig
}

func (e Event) MarshalJSON() ([]byte, error) {
	return e.libevent.MarshalJSON()
}

func (e Event) String() string {
	j, err := e.libevent.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(j)
}
