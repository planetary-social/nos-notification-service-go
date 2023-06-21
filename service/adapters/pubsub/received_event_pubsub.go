package pubsub

import (
	"context"

	"github.com/planetary-social/go-notification-service/service/domain"
)

type ReceivedEvent struct {
	Relay domain.RelayAddress
	Event domain.Event
}

type ReceivedEventPubSub struct {
	pubsub *GoChannelPubSub[ReceivedEvent]
}

func NewReceivedEventPubSub() *ReceivedEventPubSub {
	return &ReceivedEventPubSub{
		pubsub: NewGoChannelPubSub[ReceivedEvent](),
	}
}

func (m *ReceivedEventPubSub) Publish(relay domain.RelayAddress, event domain.Event) {
	m.pubsub.Publish(
		ReceivedEvent{
			Relay: relay,
			Event: event,
		},
	)
}

func (m *ReceivedEventPubSub) SubscribeToRequests(ctx context.Context) <-chan ReceivedEvent {
	return m.pubsub.Subscribe(ctx)
}
