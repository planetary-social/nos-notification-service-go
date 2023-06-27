package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/app"
	"github.com/planetary-social/go-notification-service/service/domain"
	"github.com/planetary-social/go-notification-service/service/domain/notifications"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	collectionEvents              = "events"
	collectionEventsNotifications = "notifications"
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

func (e *EventRepository) Save(event domain.Event) error {
	if err := e.saveUnderEvents(event); err != nil {
		return errors.Wrap(err, "error saving under events")
	}

	return nil
}

func (e *EventRepository) Exists(ctx context.Context, id domain.EventId) (bool, error) {
	_, err := e.client.Collection(collectionEvents).Doc(id.Hex()).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		return false, errors.Wrap(err, "error checking if document exists")
	}
	return true, nil
}

func (e *EventRepository) SaveNotificationForEvent(notification notifications.Notification) error {
	notificationDocPath := e.client.
		Collection(collectionEvents).
		Doc(notification.Event().Id().Hex()).
		Collection(collectionEventsNotifications).
		Doc(notification.UUID().String())

	notificationDocData := map[string]any{
		"uuid":    notification.UUID().String(),
		"token":   notification.APNSToken(),
		"payload": notification.Payload(),
	}

	if err := e.tx.Set(notificationDocPath, notificationDocData, firestore.MergeAll); err != nil {
		return errors.Wrap(err, "error updating the notification doc")
	}

	return nil
}

func (e *EventRepository) saveUnderEvents(event domain.Event) error {
	// todo how to handle tags? do we want to save tags in a searchable way?

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

func (e *EventRepository) GetEvents(ctx context.Context, filters domain.Filters) <-chan app.EventOrError {
	ch := make(chan app.EventOrError)
	go e.getEvents(ctx, filters, ch)
	return ch
}

func (e *EventRepository) getEvents(ctx context.Context, filters domain.Filters, ch chan<- app.EventOrError) {
	select {
	case ch <- app.NewEventOrErrorWithError(errors.New("not implemented")):
	case <-ctx.Done():
	}
}
