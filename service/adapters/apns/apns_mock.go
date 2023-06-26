package apns

import (
	"github.com/planetary-social/go-notification-service/internal/logging"
	"github.com/planetary-social/go-notification-service/service/config"
	"github.com/planetary-social/go-notification-service/service/domain/notifications"
)

type APNSMock struct {
	logger logging.Logger
}

func NewAPNSMock(config config.Config, logger logging.Logger) (*APNSMock, error) {
	return &APNSMock{logger: logger}, nil
}

func (a *APNSMock) SendNotification(notification notifications.Notification) error {
	a.logger.
		Debug().
		WithField("notification", notification.UUID()).
		Message("sending a mock APNs notification")
	return nil
}
