package di

import (
	"github.com/google/wire"
	"github.com/planetary-social/go-notification-service/service/ports/http"
	"github.com/planetary-social/go-notification-service/service/ports/pubsub"
)

var portsSet = wire.NewSet(
	http.NewServer,

	pubsub.NewReceivedEventSubscriber,
)
