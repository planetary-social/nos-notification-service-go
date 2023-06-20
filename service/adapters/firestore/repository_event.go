package firestore

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/planetary-social/go-notification-service/service/domain"
)

type EventRepository struct {
	client *firestore.Client
	tx     *firestore.Transaction
}

func NewEventRepository(client *firestore.Client, tx *firestore.Transaction) *EventRepository {
	return &EventRepository{client: client, tx: tx}
}

func (e EventRepository) Save(relay domain.RelayAddress, event domain.Event) error {
	fmt.Println("saving", relay, string(event.Content()))
	return nil
}
