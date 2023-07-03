package di

import (
	"github.com/google/wire"
	"github.com/planetary-social/go-notification-service/service/ports/firestorepubsub"
	"github.com/planetary-social/go-notification-service/service/ports/http"
	"github.com/planetary-social/go-notification-service/service/ports/memorypubsub"
)

var portsSet = wire.NewSet(
	http.NewServer,
	http.NewMetricsServer,

	memorypubsub.NewReceivedEventSubscriber,
	firestorepubsub.NewEventSavedSubscriber,
)
