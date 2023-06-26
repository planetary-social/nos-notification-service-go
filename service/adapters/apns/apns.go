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
	cfg    config.Config
	logger logging.Logger
}

func NewAPNS(cfg config.Config, logger logging.Logger) (*APNS, error) {
	client, err := newClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "error creating an apns client")
	}
	return &APNS{client: client, cfg: cfg, logger: logger}, nil
}

func newClient(cfg config.Config) (*apns2.Client, error) {
	cert, err := certificate.FromP12File(cfg.APNSCertificatePath(), cfg.APNSCertificatePassword())
	if err != nil {
		return nil, errors.Wrap(err, "error loading certificate")
	}

	switch cfg.Environment() {
	case config.EnvironmentProduction:
		return apns2.NewClient(cert).Production(), nil
	case config.EnvironmentDevelopment:
		return apns2.NewClient(cert).Development(), nil
	default:
		return nil, errors.New("unknown environment")
	}
}

func (a *APNS) SendNotification(notification notifications.Notification) error {
	n := &apns2.Notification{}
	n.ApnsID = notification.UUID().String()
	n.DeviceToken = notification.APNSToken().Hex()
	n.Topic = a.cfg.APNSTopic()
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
