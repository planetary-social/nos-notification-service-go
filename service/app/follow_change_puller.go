package app

import (
	"context"
	"time"

	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/internal/logging"
)

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
		case followChange, ok := <-ch:
			if !ok {
				return nil // Channel closed, exit gracefully
			}

			f.logger.Debug().Message(followChange.String())

			tokens, err := f.queries.GetTokens.Handle(ctx, followChange.Followee)
			if err != nil {
				// Not one of our users, ignore
				continue
			}

			for _, token := range tokens {
				if err := f.apns.SendFollowChangeNotification(*followChange, token); err != nil {
					f.logger.Error().
						WithField("token", token.Hex()).
						WithField("followee", followChange.Followee.Hex()).
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
