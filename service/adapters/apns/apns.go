package apns

import (
	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/internal/logging"
	"github.com/planetary-social/go-notification-service/service/config"
	"github.com/planetary-social/go-notification-service/service/domain/notifications"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

type APNS struct {
	client *apns2.Client
	config config.Config
	logger logging.Logger
}

func NewAPNS(config config.Config, logger logging.Logger) (*APNS, error) {
	cert, err := certificate.FromP12File(config.APNSCertificatePath(), "") // todo password support?
	if err != nil {
		return nil, errors.Wrap(err, "error loading certificate")
	}

	// If you want to test push notifications for builds running directly from XCode (Development), use
	// client := apns2.NewClient(cert).Development()
	// For apps published to the app store or installed as an ad-hoc distribution use Production()

	client := apns2.NewClient(cert).Production() // todo dev/prod

	return &APNS{client: client, config: config, logger: logger}, nil
}

func (a *APNS) SendNotification(notification notifications.Notification) error {
	n := &apns2.Notification{}
	n.ApnsID = notification.UUID().String()
	n.DeviceToken = notification.APNSToken().Hex()
	n.Topic = a.config.APNSTopic()
	n.Payload = notification.Payload()

	_, err := a.client.Push(n)
	if err != nil {
		return errors.Wrap(err, "error pushing the notification")
	}

	a.logger.Debug().
		WithField("uuid", notification.UUID().String()).
		Message("sent a notification")

	return nil
}
