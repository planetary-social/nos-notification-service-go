package di

import (
	"github.com/google/wire"
	"github.com/planetary-social/go-notification-service/service/app"
)

var downloaderSet = wire.NewSet(
	app.NewDownloader,
)
