package firestore

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill"
	watermillfirestore "github.com/ThreeDotsLabs/watermill-firestore/pkg/firestore"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/domain"
)

const PubsubTopicEventSaved = "event_saved"

type Publisher struct {
	publisher *watermillfirestore.Publisher
}

func NewPublisher(publisher *watermillfirestore.Publisher) *Publisher {
	return &Publisher{publisher: publisher}
}

func (p Publisher) PublishEventSaved(ctx context.Context, id domain.EventId) error {
	payload := EventSavedPayload{EventId: id.Hex()}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "error marshaling the payload")
	}

	msg := message.NewMessage(watermill.NewULID(), payloadJSON)
	return p.publisher.Publish(PubsubTopicEventSaved, msg)
}

type EventSavedPayload struct {
	EventId string `json:"eventId"`
}
