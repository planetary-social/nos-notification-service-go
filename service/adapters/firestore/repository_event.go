package firestore

import (
	"cloud.google.com/go/firestore"
	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/domain"
)

const (
	collectionEvents = "events"
)

type EventRepository struct {
	client          *firestore.Client
	tx              *firestore.Transaction
	relayRepository *RelayRepository
}

func NewEventRepository(
	client *firestore.Client,
	tx *firestore.Transaction,
	relayRepository *RelayRepository,
) *EventRepository {
	return &EventRepository{
		client:          client,
		tx:              tx,
		relayRepository: relayRepository,
	}
}

func (e *EventRepository) Save(relay domain.RelayAddress, event domain.Event) error {
	if err := e.saveUnderEvents(event); err != nil {
		return errors.Wrap(err, "error saving under events")
	}

	return nil
}

func (e *EventRepository) saveUnderEvents(event domain.Event) error {
	// todo how to handle tags

	eventDocPath := e.client.Collection(collectionEvents).Doc(event.Id().Hex())
	eventDocData := map[string]any{
		"id":        event.Id().Hex(),
		"publicKey": event.PubKey().Hex(),
		"createdAt": event.CreatedAt(),
		"kind":      event.Kind().Int(),
		"content":   event.Content(),
		"sig":       event.Sig().Hex(),
	}
	if err := e.tx.Set(eventDocPath, eventDocData, firestore.MergeAll); err != nil {
		return errors.Wrap(err, "error updating the event doc")
	}

	return nil
}
