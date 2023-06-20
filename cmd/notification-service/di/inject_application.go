package di

import (
	"github.com/google/wire"
	"github.com/planetary-social/go-notification-service/service/app"
)

var applicationSet = wire.NewSet(
	wire.Struct(new(app.Application), "*"),

	commandsSet,
	queriesSet,
)

var commandsSet = wire.NewSet(
	wire.Struct(new(app.Commands), "*"),

	app.NewSaveRegistrationHandler,
)

var queriesSet = wire.NewSet(
	wire.Struct(new(app.Queries), "*"),

	app.NewGetRelaysHandler,
	app.NewGetPublicKeysHandler,
)
