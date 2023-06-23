package notifications

import (
	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/domain"
	"github.com/sideshow/apns2/payload"
)

type Generator struct {
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(mention domain.PublicKey, token domain.APNSToken, event domain.Event) ([]Notification, error) {
	if mentionedThemself(mention, event) {
		return nil, nil
	}

	// todo

	return nil, nil
}

func mentionedThemself(mention domain.PublicKey, event domain.Event) bool {
	return mention == event.PubKey()
}

type Notification struct {
	token   domain.APNSToken
	payload *payload.Payload
}

func (n Notification) APNSToken() domain.APNSToken {
	return n.token
}

func (n Notification) Payload() *payload.Payload {
	return n.payload
}

type NotificationID struct {
	s string
}

func NewNotificationID(s string) (NotificationID, error) {
	if s == "" {
		return NotificationID{}, errors.New("empty id")
	}
	return NotificationID{s: s}, nil
}

func (id NotificationID) String() string {
	return id.s
}
