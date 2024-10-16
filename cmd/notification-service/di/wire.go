//go:build wireinject
// +build wireinject

package di

import (
	"context"

	googlefirestore "cloud.google.com/go/firestore"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/google/wire"
	"github.com/planetary-social/go-notification-service/service/adapters/apns"
	"github.com/planetary-social/go-notification-service/service/app"
	"github.com/planetary-social/go-notification-service/service/config"
	"github.com/planetary-social/go-notification-service/service/domain/notifications"
)

func BuildService(context.Context, config.Config) (Service, func(), error) {
	wire.Build(
		NewService,

		portsSet,
		applicationSet,
		firestoreAdaptersSet,
		downloaderSet,
		generatorSet,
		pubsubSet,
		loggingSet,
		adaptersSet,
		followChangePullerSet,
		vanishSubscriberSet,
	)
	return Service{}, nil, nil
}

type IntegrationService struct {
	Service Service

	MockAPNS *apns.APNSMock
}

func BuildIntegrationService(context.Context, config.Config) (IntegrationService, func(), error) {
	wire.Build(
		wire.Struct(new(IntegrationService), "*"),

		NewService,

		portsSet,
		applicationSet,
		firestoreAdaptersSet,
		downloaderSet,
		followChangePullerSet,
		vanishSubscriberSet,
		generatorSet,
		pubsubSet,
		loggingSet,
		integrationAdaptersSet,
	)
	return IntegrationService{}, nil, nil
}

type buildTransactionFirestoreAdaptersDependencies struct {
	LoggerAdapter watermill.LoggerAdapter
}

func buildTransactionFirestoreAdapters(client *googlefirestore.Client, tx *googlefirestore.Transaction, deps buildTransactionFirestoreAdaptersDependencies) (app.Adapters, error) {
	wire.Build(
		wire.Struct(new(app.Adapters), "*"),
		wire.FieldsOf(new(buildTransactionFirestoreAdaptersDependencies), "LoggerAdapter"),

		firestoreTxAdaptersSet,
	)
	return app.Adapters{}, nil

}

var downloaderSet = wire.NewSet(
	app.NewDownloader,
)

var followChangePullerSet = wire.NewSet(
	app.NewFollowChangePuller,
)

var vanishSubscriberSet = wire.NewSet(
	app.NewVanishSubscriber,
)

var generatorSet = wire.NewSet(
	notifications.NewGenerator,
)
