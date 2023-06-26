package main

import (
	"context"
	"fmt"
	"os"

	"github.com/boreq/errors"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/planetary-social/go-notification-service/cmd/notification-service/di"
	"github.com/planetary-social/go-notification-service/service/app"
	"github.com/planetary-social/go-notification-service/service/config"
	"github.com/planetary-social/go-notification-service/service/domain"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()
	cfg, err := config.NewConfig(
		"",
		"test-project-id",
		"some.topic",
		"./cert",
		"",
		config.EnvironmentDevelopment,
	)
	if err != nil {
		return errors.Wrap(err, "error creating a config")
	}

	service, cleanup, err := di.BuildService(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "error building a service")
	}
	defer cleanup()

	//addMyRegistration(ctx, service) // todo remove

	return service.Run(ctx)

}

func addMyRegistration(ctx context.Context, service di.Service) {
	nsec := os.Getenv("NSEC")
	_, value, err := nip19.Decode(nsec)
	if err != nil {
		panic(err)
	}

	secretKey := value.(string)
	publicKeyString, err := nostr.GetPublicKey(secretKey)
	if err != nil {
		panic(err)
	}

	publicKey, err := domain.NewPublicKeyFromHex(publicKeyString)
	if err != nil {
		panic(err)
	}

	relayAddress, err := domain.NewRelayAddress("wss://relay.damus.io")
	if err != nil {
		panic(err)
	}

	libEvent := nostr.Event{
		CreatedAt: nostr.Now(),
		Kind:      12345,
		Tags:      nostr.Tags{},
		Content: fmt.Sprintf(`
{
  "publicKey": "%s",
  "relays": [
	{
	  "address": "%s"
	}
  ],
  "apnsToken": "%s"
}
`,
			publicKey.Hex(),
			relayAddress.String(),
			"deadbeef",
		),
	}

	err = libEvent.Sign(secretKey)
	if err != nil {
		panic(err)
	}

	event, err := domain.NewEvent(libEvent)
	if err != nil {
		panic(err)
	}

	registration, err := domain.NewRegistrationFromEvent(event)
	if err != nil {
		panic(err)
	}

	cmd := app.NewSaveRegistration(registration)
	err = service.App().Commands.SaveRegistration.Handle(ctx, cmd)
	if err != nil {
		panic(err)
	}
}
