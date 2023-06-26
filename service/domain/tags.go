package domain

import (
	"github.com/boreq/errors"
)

const tagProfile = "p"

func GetMentionsFromTags(tags []EventTag) ([]PublicKey, error) {
	var mentions []PublicKey

	for _, tag := range tags {
		if tag.IsProfile() {
			pubKey, err := tag.Profile()
			if err != nil {
				return nil, errors.Wrapf(err, "error getting public key from tag '%s'", tag)
			}
			mentions = append(mentions, pubKey)
		}
	}

	return mentions, nil
}

type EventTag struct {
	tag []string
}

func NewEventTag(tag []string) (EventTag, error) {
	if len(tag) < 2 {
		return EventTag{}, errors.New("tag needs at least two fields I recon")
	}

	return EventTag{tag}, nil
}

func (e *EventTag) IsProfile() bool {
	return e.tag[0] == tagProfile
}

func (e *EventTag) Profile() (PublicKey, error) {
	if !e.IsProfile() {
		return PublicKey{}, errors.New("not a profile tag")
	}
	return NewPublicKeyFromHex(e.tag[1])
}
