package di

import (
	"context"

	googlefirestore "cloud.google.com/go/firestore"
	"github.com/boreq/errors"
	"github.com/google/wire"
	"github.com/planetary-social/go-notification-service/internal/logging"
	"github.com/planetary-social/go-notification-service/service/adapters/apns"
	"github.com/planetary-social/go-notification-service/service/adapters/firestore"
	"github.com/planetary-social/go-notification-service/service/adapters/prometheus"
	"github.com/planetary-social/go-notification-service/service/app"
	"github.com/planetary-social/go-notification-service/service/config"
)

var firestoreAdaptersSet = wire.NewSet(
	newFirestoreClient,

	firestore.NewTransactionProvider,
	wire.Bind(new(app.TransactionProvider), new(*firestore.TransactionProvider)),

	newAdaptersFactoryFn,
)

func newAdaptersFactoryFn() firestore.AdaptersFactoryFn {
	return func(client *googlefirestore.Client, tx *googlefirestore.Transaction) (app.Adapters, error) {
		return buildTransactionFirestoreAdapters(client, tx)
	}
}

var firestoreTxAdaptersSet = wire.NewSet(
	firestore.NewRegistrationRepository,
	wire.Bind(new(app.RegistrationRepository), new(*firestore.RegistrationRepository)),

	firestore.NewEventRepository,
	wire.Bind(new(app.EventRepository), new(*firestore.EventRepository)),

	firestore.NewRelayRepository,
	wire.Bind(new(app.RelayRepository), new(*firestore.RelayRepository)),

	firestore.NewPublicKeyRepository,
	wire.Bind(new(app.PublicKeyRepository), new(*firestore.PublicKeyRepository)),

	firestore.NewTagRepository,
)

var adaptersSet = wire.NewSet(
	apns.NewAPNS,
	wire.Bind(new(app.APNS), new(*apns.APNS)),

	prometheus.NewPrometheus,
	wire.Bind(new(app.Metrics), new(*prometheus.Prometheus)),
)

var integrationAdaptersSet = wire.NewSet(
	apns.NewAPNSMock,
	wire.Bind(new(app.APNS), new(*apns.APNSMock)),

	prometheus.NewPrometheus,
	wire.Bind(new(app.Metrics), new(*prometheus.Prometheus)),
)

func newFirestoreClient(ctx context.Context, config config.Config, logger logging.Logger) (*googlefirestore.Client, func(), error) {
	v, err := firestore.NewClient(ctx, config)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error creating the firestore client")
	}

	return v, func() {
		if err := v.Close(); err != nil {
			logger.Error().WithError(err).Message("error closing firestore")
		}
	}, nil
}
