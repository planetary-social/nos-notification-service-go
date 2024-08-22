package apns

import (
	"encoding/json"
	"time"

	"github.com/boreq/errors"
	"github.com/google/uuid"
	"github.com/planetary-social/go-notification-service/internal/logging"
	"github.com/planetary-social/go-notification-service/service/config"
	"github.com/planetary-social/go-notification-service/service/domain"
	"github.com/planetary-social/go-notification-service/service/domain/notifications"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

type Metrics interface {
	ReportCallToAPNS(statusCode int, err error)
}

type APNS struct {
	client  *apns2.Client
	cfg     config.Config
	metrics Metrics
	logger  logging.Logger
}

func NewAPNS(cfg config.Config, metrics Metrics, logger logging.Logger) (*APNS, error) {
	client, err := newClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "error creating an apns client")
	}
	return &APNS{
		client:  client,
		cfg:     cfg,
		metrics: metrics,
		logger:  logger.New("apns"),
	}, nil
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
	n.PushType = apns2.PushTypeBackground
	n.ApnsID = notification.UUID().String()
	n.DeviceToken = notification.APNSToken().Hex()
	n.Topic = a.cfg.APNSTopic()
	n.Payload = notification.Payload()
	n.Priority = apns2.PriorityLow

	resp, err := a.client.Push(n)
	//a.metrics.ReportCallToAPNS(resp.StatusCode, err)
	if err != nil {
		return errors.Wrap(err, "error pushing the notification")
	}

	a.logger.Debug().
		WithField("uuid", notification.UUID().String()).
		WithField("response.reason", resp.Reason).
		WithField("response.statusCode", resp.StatusCode).
		WithField("host", a.client.Host).
		Message("sent a notification")

	return nil
}

func (a *APNS) SendFollowChangeNotification(followChange domain.FollowChange, apnsToken domain.APNSToken) error {
	if apnsToken.Hex() == "" {
		return errors.New("invalid APNs token")
	}
	n, err := a.buildFollowChangeNotification(followChange, apnsToken)
	if err != nil {
		return err
	}
	resp, err := a.client.Push(n)
	//a.metrics.ReportCallToAPNS(resp.StatusCode, err)
	if err != nil {
		return errors.Wrap(err, "error pushing the follow change notification")
	}

	if resp.StatusCode == 200 {
		a.logger.Debug().
			WithField("uuid", n.ApnsID).
			WithField("response.reason", resp.Reason).
			WithField("response.statusCode", resp.StatusCode).
			WithField("host", a.client.Host).
			Message("sent a follow change notification")
	} else {
		a.logger.Error().
			WithField("uuid", n.ApnsID).
			WithField("response.reason", resp.Reason).
			WithField("response.statusCode", resp.StatusCode).
			WithField("host", a.client.Host).
			Message("failed to send a follow change notification")
	}

	return nil
}

func (a *APNS) buildFollowChangeNotification(followChange domain.FollowChange, apnsToken domain.APNSToken) (*apns2.Notification, error) {
	payload, err := followChangePayload(followChange)
	if err != nil {
		return nil, errors.Wrap(err, "error creating a payload")
	}

	n := &apns2.Notification{
		PushType:    apns2.PushTypeAlert,
		ApnsID:      uuid.New().String(),
		DeviceToken: apnsToken.Hex(),
		Topic:       a.cfg.APNSTopic(),
		Payload:     payload,
		Priority:    apns2.PriorityLow,
	}

	return n, nil
}

func followChangePayload(followChange domain.FollowChange) ([]byte, error) {
	alertMessage := ""
	if followChange.ChangeType == "unfollowed" {
		alertMessage = followChange.FriendlyFollower + " has unfollowed you!"
	} else {
		alertMessage = followChange.FriendlyFollower + " is a new follower!"
	}

	payload := map[string]interface{}{
		"aps": map[string]interface{}{
			"alert": alertMessage,
			"sound": "default",
			"badge": 1,
		},
		"data": map[string]interface{}{
			"changeType":       followChange.ChangeType,
			"at":               followChange.At.Format(time.RFC3339),
			"follower":         followChange.Follower.Hex(),
			"friendlyFollower": followChange.FriendlyFollower,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return payloadBytes, nil
}
