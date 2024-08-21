package app

import (
	"context"
	"time"

	"github.com/planetary-social/go-notification-service/internal/logging"
)

type FollowChangePuller struct {
	externalFollowChangeSubscriber ExternalFollowChangeSubscriber
	logger                         logging.Logger
	metrics                        Metrics
	counter                        int
}

func NewFollowChangePuller(
	externalFollowChangeSubscriber ExternalFollowChangeSubscriber,
	logger logging.Logger,
	metrics Metrics,
) *FollowChangePuller {
	return &FollowChangePuller{
		externalFollowChangeSubscriber: externalFollowChangeSubscriber,
		logger:                         logger.New("followChangePuller"),
		metrics:                        metrics,
		counter:                        0,
	}
}

func (f *FollowChangePuller) Run(ctx context.Context) error {
	go f.storeMetricsLoop(ctx)

	ch, err := f.externalFollowChangeSubscriber.Subscribe(ctx)
	if err != nil {
		return err
	}

	for followChange := range ch {
		f.logger.Debug().WithField("followChange", followChange).Message("received follow change")
		f.counter += 1
	}

	return nil
}

func (f *FollowChangePuller) storeMetricsLoop(ctx context.Context) {
	for {
		f.storeMetrics()

		select {
		case <-time.After(storeMetricsEvery):
		case <-ctx.Done():
			return
		}
	}
}

func (f *FollowChangePuller) storeMetrics() {
	f.metrics.MeasureFollowChange(f.counter)
	f.counter = 0
}
