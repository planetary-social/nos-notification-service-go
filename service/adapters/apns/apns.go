package apns

import (
	"fmt"
	"log"

	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/config"
	"github.com/planetary-social/go-notification-service/service/domain"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

type APNS struct {
	client *apns2.Client
}

func NewAPNS(config config.Config) (*APNS, error) {
	cert, err := certificate.FromP12File("../cert.p12", "") // todo
	if err != nil {
		return nil, errors.Wrap(err, "error loading certificate")
	}

	// If you want to test push notifications for builds running directly from XCode (Development), use
	// client := apns2.NewClient(cert).Development()
	// For apps published to the app store or installed as an ad-hoc distribution use Production()

	client := apns2.NewClient(cert).Production() // todo dev/prod

	return &APNS{client: client}, nil
}

func (a *APNS) SendNotification(token domain.APNSToken, payload []byte) error {
	notification := &apns2.Notification{}
	notification.DeviceToken = "11aa01229f15f0f0c52029d8cf8cd0aeaf2365fe4cebc4af26cd6d76b7919ef7"
	notification.Topic = "com.sideshow.Apns2"
	notification.Payload = []byte(`{"aps":{"alert":"Hello!"}}`) // See Payload section below

	res, err := a.client.Push(notification)
	if err != nil {
		return errors.Wrap(err, "error pushing a notification")
		log.Fatal("Error:", err)
	}

	fmt.Printf("%v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
}
