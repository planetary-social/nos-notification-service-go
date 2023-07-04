//go:build wireinject
// +build wireinject

package di

import (
	"context"

	googlefirestore "cloud.google.com/go/firestore"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/google/wire"
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
	)
	return Service{}, nil, nil
}

func BuildIntegrationService(context.Context, config.Config) (Service, func(), error) {
	wire.Build(
		NewService,

		portsSet,
		applicationSet,
		firestoreAdaptersSet,
		downloaderSet,
		generatorSet,
		pubsubSet,
		loggingSet,
		integrationAdaptersSet,
	)
	return Service{}, nil, nil
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

var generatorSet = wire.NewSet(
	notifications.NewGenerator,
)
