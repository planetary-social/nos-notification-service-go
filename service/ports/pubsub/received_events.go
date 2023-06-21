// Package pubsub receives internal events.
package pubsub

import (
	"context"
	"fmt"

	"github.com/planetary-social/go-notification-service/service/adapters/pubsub"
	"github.com/planetary-social/go-notification-service/service/app"
)

type ProcessReceivedEventHandler interface {
	Handle(ctx context.Context, cmd app.ProcessReceivedEvent) error
}

type ReceivedEventSubscriber struct {
	pubsub  *pubsub.ReceivedEventPubSub
	handler ProcessReceivedEventHandler
}

func NewReceivedEventSubscriber(pubsub *pubsub.ReceivedEventPubSub, handler ProcessReceivedEventHandler) *ReceivedEventSubscriber {
	return &ReceivedEventSubscriber{pubsub: pubsub, handler: handler}
}

func (p *ReceivedEventSubscriber) Run(ctx context.Context) error {
	for v := range p.pubsub.SubscribeToRequests(ctx) {
		cmd := app.NewProcessReceivedEvent(v.Relay, v.Event)
		if err := p.handler.Handle(ctx, cmd); err != nil {
			fmt.Println("error processing event", err)
		}
	}
	return nil
}
