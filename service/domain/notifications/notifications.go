package notifications

import (
	"github.com/boreq/errors"
	"github.com/google/uuid"
	"github.com/planetary-social/go-notification-service/internal/logging"
	"github.com/planetary-social/go-notification-service/service/domain"
	"github.com/sideshow/apns2/payload"
)

var (
	eventKindNote                   = domain.MustNewEventKind(1)
	eventKindReaction               = domain.MustNewEventKind(7)
	eventKindEncryptedDirectMessage = domain.MustNewEventKind(4)
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

	notification, err := NewNotification(event, id, token, payloadJSON)
	if err != nil {
		return nil, errors.Wrap(err, "error creating a notification")
	}

	return []Notification{notification}, nil
}

func (g *Generator) createPayload(mention domain.PublicKey, event domain.Event) ([]byte, error) {
	payload, err := g.generatePayload(mention, event)
	if err != nil {
		return nil, errors.Wrap(err, "error creating the payload")
	}

	if payload == nil {
		return nil, nil
	}

	eventJSON, err := event.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling event")
	}

	payload = payload.Custom("event", string(eventJSON))

	payloadJSON, err := payload.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling payload")
	}

	return payloadJSON, nil
}

func (g *Generator) generatePayload(mention domain.PublicKey, event domain.Event) (*payload.Payload, error) {
	if mentionedThemself(mention, event) {
		return nil, nil
	}

	switch event.Kind() {
	case eventKindNote:
		g.logger.Debug().Message("note")
		return payload.NewPayload().Alert("You were tagged in a note.").Category("event.tagged.note"), nil
	case eventKindReaction:
		g.logger.Debug().Message("reaction")
		return payload.NewPayload().Alert("You were tagged in a reaction.").Category("event.tagged.reaction"), nil
	case eventKindEncryptedDirectMessage:
		g.logger.Debug().Message("encrypted direct message")
		return payload.NewPayload().Alert("You were tagged in an encrypted direct message.").Category("event.tagged.encryptedDirectMessage"), nil
	default:
		return nil, nil
	}
}

func mentionedThemself(mention domain.PublicKey, event domain.Event) bool {
	return mention == event.PubKey()
}

type Notification struct {
	event domain.Event

	uuid    NotificationUUID
	token   domain.APNSToken
	payload []byte
}

func NewNotification(
	event domain.Event,
	uuid NotificationUUID,
	token domain.APNSToken,
	payload []byte,
) (Notification, error) {
	if len(payload) == 0 {
		return Notification{}, errors.New("empty payload")
	}
	return Notification{
		event:   event,
		uuid:    uuid,
		token:   token,
		payload: payload,
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
