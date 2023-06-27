package domain

import (
	"time"

	"github.com/boreq/errors"
	"github.com/nbd-wtf/go-nostr"
)

type Filter struct {
	ids     []EventId
	kinds   []EventKind
	authors []PublicKey
	tags    []EventTag
	since   *time.Time
	until   *time.Time
	limit   int
	search  string

	libfilter nostr.Filter
}

func NewFilter(f nostr.Filter) (Filter, error) {
	var ids []EventId
	for _, v := range f.IDs {
		id, err := NewEventId(v)
		if err != nil {
			return Filter{}, errors.Wrap(err, "error creating an event id")
		}
		ids = append(ids, id)
	}

	var kinds []EventKind
	for _, v := range f.Kinds {
		kind, err := NewEventKind(v)
		if err != nil {
			return Filter{}, errors.Wrap(err, "error creating an event kind")
		}
		kinds = append(kinds, kind)
	}

	var authors []PublicKey
	for _, v := range f.Authors {
		author, err := NewPublicKeyFromHex(v)
		if err != nil {
			return Filter{}, errors.Wrap(err, "error creating a public key")
		}
		authors = append(authors, author)
	}

	var tags []EventTag
	for _, v := range f.Tags {
		tag, err := NewEventTag(v)
		if err != nil {
			return Filter{}, errors.Wrap(err, "error creating a tag")
		}
		tags = append(tags, tag)
	}

	var since *time.Time
	if f.Since != nil {
		t := f.Since.Time()
		since = &t
	}

	var until *time.Time
	if f.Until != nil {
		t := f.Until.Time()
		until = &t
	}

	if f.Limit < 0 {
		return Filter{}, errors.New("limit can't be negative")
	}

	return Filter{
		ids:     ids,
		kinds:   kinds,
		authors: authors,
		tags:    tags,
		since:   since,
		until:   until,
		limit:   f.Limit,
		search:  f.Search,

		libfilter: f,
	}, nil
}

func (f Filter) Matches(event Event) bool {
	return f.libfilter.Matches(&event.libevent)
}

type Filters struct {
	filters []Filter
}

func NewFilters(f nostr.Filters) (Filters, error) {
	var filters []Filter
	for _, v := range f {
		filter, err := NewFilter(v)
		if err != nil {
			return Filters{}, errors.Wrap(err, "error creating a filter")
		}
		filters = append(filters, filter)
	}

	return Filters{
		filters: filters,
	}, nil
}

func (f Filters) Match(event Event) bool {
	for _, filter := range f.filters {
		if filter.Matches(event) {
			return true
		}
	}
	return false
}
