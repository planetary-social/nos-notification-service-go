package pubsub

import (
	"context"

	"github.com/planetary-social/go-notification-service/service/app"
	"github.com/planetary-social/go-notification-service/service/domain"
)

type ReceivedEventPubSub struct {
	pubsub *GoChannelPubSub[app.ReceivedEvent]
}

func NewReceivedEventPubSub() *ReceivedEventPubSub {
	return &ReceivedEventPubSub{
		pubsub: NewGoChannelPubSub[app.ReceivedEvent](),
	}
}

func (m *ReceivedEventPubSub) Publish(relay domain.RelayAddress, event domain.Event) {
	m.pubsub.Publish(
		app.NewReceivedEvent(relay, event),
	)
}

func (m *ReceivedEventPubSub) Subscribe(ctx context.Context) <-chan app.ReceivedEvent {
	return m.pubsub.Subscribe(ctx)
}
