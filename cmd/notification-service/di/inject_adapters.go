package di

import (
	"github.com/google/wire"
	"github.com/planetary-social/go-notification-service/service/adapters/firestore"
	"github.com/planetary-social/go-notification-service/service/app"
)

var firestoreAdaptersSet = wire.NewSet(
	firestore.NewTransactionProvider,
	wire.Bind(new(app.TransactionProvider), new(*firestore.TransactionProvider)),
)

var adaptersSet = wire.NewSet(
// adapters.NewCurrentTimeProvider,
// wire.Bind(new(commands.CurrentTimeProvider), new(*adapters.CurrentTimeProvider)),
// wire.Bind(new(boxstream.CurrentTimeProvider), new(*adapters.CurrentTimeProvider)),
// wire.Bind(new(invitesadapters.CurrentTimeProvider), new(*adapters.CurrentTimeProvider)),
// wire.Bind(new(blobreplication.CurrentTimeProvider), new(*adapters.CurrentTimeProvider)),
)
