package notifications

import (
	"time"

	"github.com/boreq/errors"
	"github.com/google/uuid"
	"github.com/planetary-social/go-notification-service/internal/logging"
	"github.com/planetary-social/go-notification-service/service/domain"
	"github.com/sideshow/apns2/payload"
)

type Generator struct {
	logger logging.Logger
}

func NewGenerator(logger logging.Logger) *Generator {
	return &Generator{
		logger: logger.New("generator"),
	}
}

func (g *Generator) Generate(mention domain.PublicKey, token domain.APNSToken, event domain.Event) ([]Notification, error) {
	payloadJSON, err := g.createPayload(mention, event)
	if err != nil {
		return nil, errors.Wrap(err, "error creating the payload")
	}

	if payloadJSON == nil {
		return nil, nil
	}

	id, err := NewNotificationUUID()
	if err != nil {
		return nil, errors.Wrap(err, "error generating a notification id")
	}

	notification, err := NewNotification(event, id, token, payloadJSON, time.Now())
	if err != nil {
		return nil, errors.Wrap(err, "error creating a notification")
	}

	return []Notification{notification}, nil
}

func (g *Generator) createPayload(mention domain.PublicKey, event domain.Event) ([]byte, error) {
	if g.mentionedThemself(mention, event) {
		return nil, nil
	}

	notificationPayload := payload.NewPayload().ContentAvailable()

	payloadJSON, err := notificationPayload.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling payload")
	}

	return payloadJSON, nil
}

func (g *Generator) mentionedThemself(mention domain.PublicKey, event domain.Event) bool {
	return mention == event.PubKey()
}

type Notification struct {
	event domain.Event

	uuid      NotificationUUID
	token     domain.APNSToken
	payload   []byte
	createdAt *time.Time // old notifications don't have this value
}

func NewNotification(
	event domain.Event,
	uuid NotificationUUID,
	token domain.APNSToken,
	payload []byte,
	createdAt time.Time,
) (Notification, error) {
	if len(payload) == 0 {
		return Notification{}, errors.New("empty payload")
	}
	return Notification{
		event:     event,
		uuid:      uuid,
		token:     token,
		payload:   payload,
		createdAt: &createdAt,
	}, nil
}

func NewNotificationFromHistory(
	event domain.Event,
	uuid NotificationUUID,
	token domain.APNSToken,
	payload []byte,
	createdAt *time.Time,
) (Notification, error) {
	if len(payload) == 0 {
		return Notification{}, errors.New("empty payload")
	}
	return Notification{
		event:     event,
		uuid:      uuid,
		token:     token,
		payload:   payload,
		createdAt: createdAt,
	}, nil
}

func (n Notification) Event() domain.Event {
	return n.event
}

func (n Notification) UUID() NotificationUUID {
	return n.uuid
}

func (n Notification) APNSToken() domain.APNSToken {
	return n.token
}

func (n Notification) Payload() []byte {
	return n.payload
}

func (n Notification) CreatedAt() *time.Time {
	return n.createdAt
}

type NotificationUUID struct {
	s string
}

func NewNotificationUUID() (NotificationUUID, error) {
	return NewNotificationUUIDFromString(uuid.New().String())
}

func NewNotificationUUIDFromString(s string) (NotificationUUID, error) {
	if s == "" {
		return NotificationUUID{}, errors.New("empty id")
	}
	_, err := uuid.Parse(s)
	if err != nil {
		return NotificationUUID{}, errors.Wrap(err, "malformed uuid")
	}
	return NotificationUUID{s: s}, nil
}

func (id NotificationUUID) String() string {
	return id.s
}
