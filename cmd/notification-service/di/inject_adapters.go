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
	return func(transaction *googlefirestore.Transaction) (app.Adapters, error) {
		return buildTransactionFirestoreAdapters(transaction)
	}
}

var firestoreTxAdaptersSet = wire.NewSet(
	firestore.NewRegistrationRepository,
	wire.Bind(new(app.RegistrationRepository), new(*firestore.RegistrationRepository)),
)

var adaptersSet = wire.NewSet(
// adapters.NewCurrentTimeProvider,
// wire.Bind(new(commands.CurrentTimeProvider), new(*adapters.CurrentTimeProvider)),
// wire.Bind(new(boxstream.CurrentTimeProvider), new(*adapters.CurrentTimeProvider)),
// wire.Bind(new(invitesadapters.CurrentTimeProvider), new(*adapters.CurrentTimeProvider)),
// wire.Bind(new(blobreplication.CurrentTimeProvider), new(*adapters.CurrentTimeProvider)),
)
