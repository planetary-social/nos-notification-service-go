//go:build wireinject
// +build wireinject

package di

import (
	"context"

	googlefirestore "cloud.google.com/go/firestore"
	"github.com/google/wire"
	"github.com/planetary-social/go-notification-service/service/app"
	"github.com/planetary-social/go-notification-service/service/config"
)

func BuildService(context.Context, config.Config) (Service, func(), error) {
	wire.Build(
		NewService,

		portsSet,
		applicationSet,
		firestoreAdaptersSet,
	)
	return Service{}, nil, nil
}

func buildTransactionFirestoreAdapters(tx *googlefirestore.Transaction) (app.Adapters, error) {
	wire.Build(
		wire.Struct(new(app.Adapters), "*"),

		firestoreTxAdaptersSet,
	)
	return app.Adapters{}, nil

}

//func newAdvertiser(l identity.Public, config service.Config) (*local.Advertiser, error) {
//	return local.NewAdvertiser(l, config.ListenAddress)
//}
//
//func newIntegrationTestConfig(t *testing.T) service.Config {
//	dataDirectory := fixtures.Directory(t)
//	oldDataDirectory := fixtures.Directory(t)
//
//	cfg := service.Config{
//		DataDirectory:      dataDirectory,
//		GoSSBDataDirectory: oldDataDirectory,
//		NetworkKey:         fixtures.SomeNetworkKey(),
//		MessageHMAC:        fixtures.SomeMessageHMAC(),
//	}
//	cfg.SetDefaults()
//	return cfg
//}
//
//func newBadger(system logging.LoggingSystem, logger logging.Logger, config service.Config) (*badger.DB, func(), error) {
//	badgerDirectory := filepath.Join(config.DataDirectory, "badger")
//
//	options := badger.DefaultOptions(badgerDirectory)
//	options.Logger = badgeradapters.NewLogger(system, badgeradapters.LoggerLevelWarning)
//
//	if config.ModifyBadgerOptions != nil {
//		adapter := service.NewBadgerOptionsAdapter(&options)
//		config.ModifyBadgerOptions(adapter)
//	}
//
//	db, err := badger.Open(options)
//	if err != nil {
//		return nil, nil, errors.Wrap(err, "failed to open the database")
//	}
//
//	return db, func() {
//		if err := db.Close(); err != nil {
//			logger.Error().WithError(err).Message("error closing the database")
//		}
//	}, nil
//}
//
//func privateIdentityToPublicIdentity(p identity.Private) identity.Public {
//	return p.Public()
//}
//
//func newContextLogger(loggingSystem logging.LoggingSystem) logging.Logger {
//	return logging.NewContextLogger(loggingSystem, "scuttlego")
//}
