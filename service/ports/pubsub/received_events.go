// Package pubsub receives internal events.
package pubsub

import (
	"context"

	"github.com/planetary-social/go-notification-service/internal/logging"
	"github.com/planetary-social/go-notification-service/service/adapters/pubsub"
	"github.com/planetary-social/go-notification-service/service/app"
)

type ProcessReceivedEventHandler interface {
	Handle(ctx context.Context, cmd app.ProcessReceivedEvent) error
}

type ReceivedEventSubscriber struct {
	pubsub  *pubsub.ReceivedEventPubSub
	handler ProcessReceivedEventHandler
	logger  logging.Logger
}

func NewReceivedEventSubscriber(
	pubsub *pubsub.ReceivedEventPubSub,
	handler ProcessReceivedEventHandler,
	logger logging.Logger,
) *ReceivedEventSubscriber {
	return &ReceivedEventSubscriber{
		pubsub:  pubsub,
		handler: handler,
		logger:  logger,
	}
}

func (p *ReceivedEventSubscriber) Run(ctx context.Context) error {
	for v := range p.pubsub.Subscribe(ctx) {
		cmd := app.NewProcessReceivedEvent(v.Relay, v.Event)
		if err := p.handler.Handle(ctx, cmd); err != nil {
			p.logger.Error().
				WithError(err).
				WithField("relay", v.Relay).
				WithField("event", v.Event).
				Message("error handling a received event")
		}
	}
	return nil
}
