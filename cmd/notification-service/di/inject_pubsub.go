package di

import (
	"github.com/google/wire"
	"github.com/planetary-social/go-notification-service/service/adapters/gcp"
	"github.com/planetary-social/go-notification-service/service/adapters/pubsub"
	"github.com/planetary-social/go-notification-service/service/app"
)

var pubsubSet = wire.NewSet(
	pubsub.NewReceivedEventPubSub,
	wire.Bind(new(app.ReceivedEventPublisher), new(*pubsub.ReceivedEventPubSub)),
	wire.Bind(new(app.ReceivedEventSubscriber), new(*pubsub.ReceivedEventPubSub)),
)

var googlePubsubSet = wire.NewSet(
	gcp.NewWatermillPublisher,
	gcp.NewPublisher,
	wire.Bind(new(app.ExternalEventPublisher), new(*gcp.Publisher)),
)
