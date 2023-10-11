package gcp

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/planetary-social/go-notification-service/service/config"
	"github.com/planetary-social/go-notification-service/service/domain"
	"google.golang.org/api/option"
)

const googlePubSubNostrEventsTopic = "nostr-events"

func NewWatermillPublisher(config config.Config, logger watermill.LoggerAdapter) (*googlecloud.Publisher, error) {
	var options []option.ClientOption

	if j := config.GooglePubSubCredentialsJSON(); len(j) > 0 {
		options = append(options, option.WithCredentialsJSON(config.GooglePubSubCredentialsJSON()))
	}

	publisherConfig := googlecloud.PublisherConfig{
		ProjectID:                 config.GooglePubSubProjectID(),
		DoNotCreateTopicIfMissing: true,
		ClientOptions:             options,
	}

	return googlecloud.NewPublisher(publisherConfig, logger)
}

type Publisher struct {
	publisher *googlecloud.Publisher
}

func NewPublisher(publisher *googlecloud.Publisher) *Publisher {
	return &Publisher{publisher: publisher}
}

func (p *Publisher) PublishNewEventReceived(ctx context.Context, event domain.Event) error {
	msg := message.NewMessage(watermill.NewULID(), event.Raw())
	return p.publisher.Publish(googlePubSubNostrEventsTopic, msg)
}
