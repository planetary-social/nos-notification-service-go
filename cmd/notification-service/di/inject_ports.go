package di

import (
	"github.com/google/wire"
	"github.com/planetary-social/go-notification-service/service/ports/http"
)

var portsSet = wire.NewSet(
	http.NewServer,
)
