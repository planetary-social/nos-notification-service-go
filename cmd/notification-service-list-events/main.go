package main

import (
	"context"
	"fmt"
	"os"

	"github.com/boreq/errors"
	"github.com/nbd-wtf/go-nostr"
	"github.com/planetary-social/go-notification-service/cmd/notification-service/di"
	configadapters "github.com/planetary-social/go-notification-service/service/adapters/config"
	"github.com/planetary-social/go-notification-service/service/domain"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	if len(os.Args) != 2 {
		return errors.New("usage: program <npub>")
	}

	publicKey, err := domain.NewPublicKeyFromNpub(os.Args[1])
	if err != nil {
		return errors.Wrap(err, "error decoding the npub")
	}

	cfg, err := configadapters.NewEnvironmentConfigLoader().Load()
	if err != nil {
		return errors.Wrap(err, "error creating a config")
	}

	service, cleanup, err := di.BuildService(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "error building a service")
	}
	defer cleanup()

	if err := listTokens(ctx, service, publicKey); err != nil {
		return errors.Wrap(err, "error listing APNs tokens")
	}

	if err := listEvents(ctx, service, publicKey); err != nil {
		return errors.Wrap(err, "error listing APNs tokens")
	}

	if err := listRelays(ctx, service, publicKey); err != nil {
		return errors.Wrap(err, "error listing relays")
	}

	return nil
}

func listTokens(ctx context.Context, service di.Service, publicKey domain.PublicKey) error {
	fmt.Println()
	fmt.Println("listing stored APNs tokens related to public key", publicKey.Hex())

	tokens, err := service.App().Queries.GetTokens.Handle(ctx, publicKey)
	if err != nil {
		return errors.Wrap(err, "error getting APNs tokens")
	}

	if len(tokens) == 0 {
		fmt.Println("no stored tokens")
	}

	for _, token := range tokens {
		fmt.Println(token.Hex())
	}

	return nil
}

func listRelays(ctx context.Context, service di.Service, publicKey domain.PublicKey) error {
	fmt.Println()
	fmt.Println("listing relays related to public key", publicKey.Hex())

	relays, err := service.App().Queries.GetRelays.Handle(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting relays")
	}

	for _, relay := range relays {
		publicKeys, err := service.App().Queries.GetPublicKeys.Handle(ctx, relay)
		if err != nil {
			return errors.Wrap(err, "error getting relays")
		}

		for _, publicKeyForRelay := range publicKeys {
			if publicKeyForRelay == publicKey {
				fmt.Println(relay)
				break
			}
		}
	}

	return nil
}

func listEvents(ctx context.Context, service di.Service, publicKey domain.PublicKey) error {
	fmt.Println()
	fmt.Println("listing stored events related to public key", publicKey.Hex())

	filters, err := domain.NewFilters(nostr.Filters{
		{
			Authors: []string{
				publicKey.Hex(),
			},
			Limit: 100,
		},
	})
	if err != nil {
		return errors.Wrap(err, "error creating filters")
	}

	for v := range service.App().Queries.GetEvents.Handle(ctx, filters) {
		if err := v.Err(); err != nil {
			return errors.Wrap(err, "handler returned an error")
		}

		if v.EOSE() {
			fmt.Println("end of stored events")
			break
		}

		evt := v.Event()

		notifications, err := service.App().Queries.GetNotifications.Handle(ctx, v.Event().Id())
		if err != nil {
			return errors.Wrapf(err, "error getting notifications for event '%s'", v.Event().Id().Hex())
		}

		fmt.Println("event", evt.Id().Hex(), evt.Kind().Int())

		for _, notification := range notifications {
			fmt.Println("notification", notification.UUID(), "created at", notification.CreatedAt())
		}
	}

	return nil
}
