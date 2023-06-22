package di

import (
	"github.com/google/wire"
	"github.com/planetary-social/go-notification-service/internal/logging"
	"github.com/sirupsen/logrus"
)

var loggingSet = wire.NewSet(
	newSystemLogger,

	logging.NewLogrusLoggingSystem,
	wire.Bind(new(logging.LoggingSystem), new(logging.LogrusLoggingSystem)),

	newLogrus,
)

func newSystemLogger(system logging.LoggingSystem) logging.Logger {
	return logging.NewSystemLogger(system, "root")
}

func newLogrus() *logrus.Logger {
	v := logrus.New()
	v.SetLevel(logrus.DebugLevel)
	return v
}
