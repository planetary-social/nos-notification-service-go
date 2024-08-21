package gcp

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/pkg/errors"
	"github.com/planetary-social/go-notification-service/service/config"
	"github.com/planetary-social/go-notification-service/service/domain"
	"google.golang.org/api/option"
)

const googlePubSubFollowChangeTopic = "follow-changes"

func NewWatermillSubscriber(config config.Config, logger watermill.LoggerAdapter) (*googlecloud.Subscriber, error) {
	var options []option.ClientOption

	if j := config.GooglePubSubCredentialsJSON(); len(j) > 0 {
		options = append(options, option.WithCredentialsJSON(config.GooglePubSubCredentialsJSON()))
	}

	publisherConfig := googlecloud.SubscriberConfig{
		ProjectID:                 config.GooglePubSubProjectID(),
		DoNotCreateTopicIfMissing: true,
		ClientOptions:             options,
	}

	return googlecloud.NewSubscriber(publisherConfig, logger)
}

type GCPFollowChangeSubscriber struct {
	subscriber *googlecloud.Subscriber
	logger     watermill.LoggerAdapter
}

func NewFollowChangeSubscriber(subscriber *googlecloud.Subscriber, logger watermill.LoggerAdapter) *GCPFollowChangeSubscriber {
	return &GCPFollowChangeSubscriber{subscriber: subscriber, logger: logger}
}

func (p *GCPFollowChangeSubscriber) Subscribe(ctx context.Context) (<-chan *domain.FollowChange, error) {
	subChan, err := p.subscriber.Subscribe(ctx, googlePubSubFollowChangeTopic)
	if err != nil {
		return nil, errors.Wrap(err, "error subscribing")
	}

	ch := make(chan *domain.FollowChange)

	go func() {
		defer close(ch)
		defer p.subscriber.Close()

		for message := range subChan {
			// We never retry messages so we can ACK immediately.
			message.Ack()

			var payload domain.FollowChange
			if err := json.Unmarshal(message.Payload, &payload); err != nil {
				p.logger.Error("error unmarshaling follow change payload", err, watermill.LogFields{"payload": string(message.Payload)})
				continue
			}

			select {
			case ch <- &payload:
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch, nil
}
