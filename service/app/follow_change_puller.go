package app

import (
	"context"
	"time"

	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/internal/logging"
)

// Reads from the follow-change puller, creates FollowChangeBatch types from
// each entry and sends notifications to those users for which we have APNS
// tokens
type FollowChangePuller struct {
	externalFollowChangeSubscriber ExternalFollowChangeSubscriber
	apns                           APNS
	queries                        Queries
	logger                         logging.Logger
	metrics                        Metrics
	counter                        int
}

func NewFollowChangePuller(
	externalFollowChangeSubscriber ExternalFollowChangeSubscriber,
	apns APNS,
	queries Queries,
	logger logging.Logger,
	metrics Metrics,
) *FollowChangePuller {
	return &FollowChangePuller{
		externalFollowChangeSubscriber: externalFollowChangeSubscriber,
		apns:                           apns,
		queries:                        queries,
		logger:                         logger.New("followChangePuller"),
		metrics:                        metrics,
		counter:                        0,
	}
}

// Listens for messages from the follow-change pubsub and for each of them, if
// they belong to one of our users, we send a notification
func (f *FollowChangePuller) Run(ctx context.Context) error {
	go f.storeMetricsLoop(ctx)

	ch, err := f.externalFollowChangeSubscriber.Subscribe(ctx)
	if err != nil {
		return errors.Wrap(err, "error subscribing to follow changes")
	}

	for {
		select {
		case followChangeAggregate, ok := <-ch:
			if !ok {
				return nil // Channel closed, exit gracefully
			}

			tokens, err := f.queries.GetTokens.Handle(ctx, followChangeAggregate.Followee)
			if err != nil {
				f.logger.Error().
					WithField("followee", followChangeAggregate.Followee.Hex()).
					WithError(err).
					Message("error getting tokens for followee")
				continue
			}

			if len(tokens) == 0 {
				// Not one of our users, ignore
				continue
			}

			f.logger.Debug().Message(followChangeAggregate.String())

			for _, token := range tokens {
				if err := f.apns.SendFollowChangeNotification(*followChangeAggregate, token); err != nil {
					f.logger.Error().
						WithField("token", token.Hex()).
						WithField("followee", followChangeAggregate.Followee.Hex()).
						WithError(err).
						Message("error sending follow change notification")
					continue
				}
			}

			f.counter += 1
		case <-ctx.Done():
			f.logger.Debug().Message("context canceled, shutting down FollowChangePuller")
			return nil
		}
	}
}

func (f *FollowChangePuller) storeMetricsLoop(ctx context.Context) {
	ticker := time.NewTicker(storeMetricsEvery)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			f.storeMetrics()
		case <-ctx.Done():
			f.logger.Debug().Message("context canceled, stopping metrics loop")
			return
		}
	}
}

func (f *FollowChangePuller) storeMetrics() {
	f.metrics.MeasureFollowChange(f.counter)
	f.counter = 0
}
