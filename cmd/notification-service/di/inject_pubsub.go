package di

import (
	"github.com/google/wire"
	"github.com/planetary-social/go-notification-service/service/adapters/pubsub"
	"github.com/planetary-social/go-notification-service/service/app"
)

var pubsubSet = wire.NewSet(
	pubsub.NewReceivedEventPubSub,
	wire.Bind(new(app.ReceivedEventPublisher), new(*pubsub.ReceivedEventPubSub)),
)
