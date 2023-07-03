package firestorepubsub

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/internal/logging"
	"github.com/planetary-social/go-notification-service/service/adapters/firestore"
	"github.com/planetary-social/go-notification-service/service/app"
	"github.com/planetary-social/go-notification-service/service/domain"
)

type ProcessSavedEventHandler interface {
	Handle(ctx context.Context, cmd app.ProcessSavedEvent) error
}

type FirestoreSubscriber interface {
	message.Subscriber
}

type EventSavedSubscriber struct {
	subscriber FirestoreSubscriber
	handler    ProcessSavedEventHandler
	logger     logging.Logger
}

func NewEventSavedSubscriber(
	subscriber FirestoreSubscriber,
	handler ProcessSavedEventHandler,
	logger logging.Logger,
) *EventSavedSubscriber {
	return &EventSavedSubscriber{
		subscriber: subscriber,
		handler:    handler,
		logger:     logger.New("eventSavedSubscriber"),
	}
}

func (p *EventSavedSubscriber) Run(ctx context.Context) error {
	ch, err := p.subscriber.Subscribe(ctx, firestore.PubsubTopicEventSaved)
	if err != nil {
		return errors.Wrap(err, "error subscribing")
	}

	for msg := range ch {
		if err := p.handleMessage(ctx, msg); err != nil {
			msg.Nack()
			p.logger.Error().WithError(err).Message("error handling a message")
			continue
		}

		msg.Ack()
	}
	return nil
}

func (p *EventSavedSubscriber) handleMessage(ctx context.Context, msg *message.Message) error {
	var payload firestore.EventSavedPayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return errors.Wrap(err, "error unmarshaling event payload")
	}

	eventId, err := domain.NewEventId(payload.EventId)
	if err != nil {
		return errors.Wrap(err, "error creating event id")
	}

	cmd := app.NewProcessSavedEvent(eventId)
	if err := p.handler.Handle(ctx, cmd); err != nil {
		return errors.Wrap(err, "error calling the handler")
	}

	return nil
}
