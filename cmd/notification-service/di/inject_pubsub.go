package di

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/boreq/errors"
	"github.com/google/wire"
	"github.com/planetary-social/go-notification-service/service/adapters/gcp"
	"github.com/planetary-social/go-notification-service/service/adapters/pubsub"
	"github.com/planetary-social/go-notification-service/service/app"
	"github.com/planetary-social/go-notification-service/service/config"
)

var pubsubSet = wire.NewSet(
	pubsub.NewReceivedEventPubSub,
	wire.Bind(new(app.ReceivedEventPublisher), new(*pubsub.ReceivedEventPubSub)),
	wire.Bind(new(app.ReceivedEventSubscriber), new(*pubsub.ReceivedEventPubSub)),
)

var googlePubsubSet = wire.NewSet(
	newExternalEventPublisher,
	newExternalFollowChangeSubscriber,
)

func newExternalEventPublisher(config config.Config, logger watermill.LoggerAdapter) (app.ExternalEventPublisher, error) {
	if config.GooglePubSubEnabled() {
		publisher, err := gcp.NewWatermillPublisher(config, logger)
		if err != nil {
			return nil, errors.Wrap(err, "error creating a watermil publisher")
		}
		return gcp.NewPublisher(publisher), nil
	} else {
		return gcp.NewNoopPublisher(), nil
	}
}

func newExternalFollowChangeSubscriber(config config.Config, logger watermill.LoggerAdapter) (app.ExternalFollowChangeSubscriber, error) {
	if config.GooglePubSubEnabled() {
		subscriber, err := gcp.NewWatermillSubscriber(config, logger)
		if err != nil {
			return nil, errors.Wrap(err, "error creating a watermil subscriber")
		}
		return gcp.NewFollowChangeSubscriber(subscriber, logger), nil
	} else {
		return gcp.NewNoopSubscriber(), nil
	}
}
