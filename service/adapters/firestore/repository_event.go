package firestore

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/domain"
)

const (
	collectionEvents = "events"
)

type EventRepository struct {
	client *firestore.Client
	tx     *firestore.Transaction
}

func NewEventRepository(client *firestore.Client, tx *firestore.Transaction) *EventRepository {
	return &EventRepository{client: client, tx: tx}
}

func (e *EventRepository) Save(relay domain.RelayAddress, event domain.Event) error {
	tokenDocPath := e.client.Collection(collectionEvents).Doc(event.Id().Hex())
	tokenDocData := map[string]any{
		"id":        event.Id().Hex(),
		"publicKey": event.PubKey().Hex(),
		"createdAt": event.CreatedAt(),
		"kind":      event.Kind().Int(),
		//"tags":
		"content": event.Content(),
		"sig":     event.Sig().Hex(),
	}
	if err := r.tx.Set(tokenDocPath, tokenDocData, firestore.MergeAll); err != nil {
		return errors.Wrap(err, "error updating the token doc")
	}

	fmt.Println("saving", relay, string(event.Content()))
	return nil
}
