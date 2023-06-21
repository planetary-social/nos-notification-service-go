package di

import (
	"github.com/google/wire"
	"github.com/planetary-social/go-notification-service/service/app"
	"github.com/planetary-social/go-notification-service/service/ports/pubsub"
)

var applicationSet = wire.NewSet(
	wire.Struct(new(app.Application), "*"),

	commandsSet,
	queriesSet,
)

var commandsSet = wire.NewSet(
	wire.Struct(new(app.Commands), "*"),

	app.NewSaveRegistrationHandler,

	app.NewProcessReceivedEventHandler,
	wire.Bind(new(pubsub.ProcessReceivedEventHandler), new(*app.ProcessReceivedEventHandler)),
)

var queriesSet = wire.NewSet(
	wire.Struct(new(app.Queries), "*"),

	app.NewGetRelaysHandler,
	app.NewGetPublicKeysHandler,
)
