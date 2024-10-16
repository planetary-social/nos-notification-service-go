package apns

import (
	"sync"

	"github.com/planetary-social/go-notification-service/internal"
	"github.com/planetary-social/go-notification-service/internal/logging"
	"github.com/planetary-social/go-notification-service/service/config"
	"github.com/planetary-social/go-notification-service/service/domain"
	"github.com/planetary-social/go-notification-service/service/domain/notifications"
)

type APNSMock struct {
	logger logging.Logger

	sentNotificationsLock sync.Mutex
	sentNotifications     []notifications.Notification
}

func NewAPNSMock(config config.Config, logger logging.Logger) (*APNSMock, error) {
	return &APNSMock{logger: logger}, nil
}

func (a *APNSMock) SendNotification(notification notifications.Notification) error {
	a.sentNotificationsLock.Lock()
	defer a.sentNotificationsLock.Unlock()

	a.sentNotifications = append(a.sentNotifications, notification)

	a.logger.
		Debug().
		WithField("notification", notification.UUID()).
		Message("sending a mock APNs notification")

	return nil
}

func (a *APNSMock) SendFollowChangeNotification(followChange domain.FollowChangeBatch, token domain.APNSToken) error {
	notification := notifications.Notification{}

	return a.SendNotification(notification)
}

func (a *APNSMock) SendSilentFollowChangeNotification(followChange domain.FollowChangeBatch, token domain.APNSToken) error {
	notification := notifications.Notification{}

	return a.SendNotification(notification)
}

func (a *APNSMock) SentNotifications() []notifications.Notification {
	a.sentNotificationsLock.Lock()
	defer a.sentNotificationsLock.Unlock()

	return internal.CopySlice(a.sentNotifications)
}
