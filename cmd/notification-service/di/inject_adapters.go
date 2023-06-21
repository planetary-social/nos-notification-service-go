package di

import (
	googlefirestore "cloud.google.com/go/firestore"
	"github.com/google/wire"
	"github.com/planetary-social/go-notification-service/service/adapters/firestore"
	"github.com/planetary-social/go-notification-service/service/app"
)

var firestoreAdaptersSet = wire.NewSet(
	firestore.NewClient,

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
)
