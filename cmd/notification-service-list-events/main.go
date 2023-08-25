package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/boreq/errors"
	"github.com/nbd-wtf/go-nostr"
	"github.com/planetary-social/go-notification-service/cmd/notification-service/di"
	"github.com/planetary-social/go-notification-service/internal"
	configadapters "github.com/planetary-social/go-notification-service/service/adapters/config"
	"github.com/planetary-social/go-notification-service/service/domain"
	"github.com/planetary-social/go-notification-service/service/domain/notifications"
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

	tokensSet := internal.NewEmptySet[domain.APNSToken]()

	if err := listTokens(ctx, service, publicKey, tokensSet); err != nil {
		return errors.Wrap(err, "error listing APNs tokens")
	}

	if err := listEvents(ctx, service, publicKey, tokensSet); err != nil {
		return errors.Wrap(err, "error listing APNs tokens")
	}

	return nil
}

func listTokens(ctx context.Context, service di.Service, publicKey domain.PublicKey, tokensSet *internal.Set[domain.APNSToken]) error {
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
		tokensSet.Put(token)
		fmt.Println("token", token.Hex())
	}

	return nil
}

type eventWithNotifications struct {
	Event         domain.Event
	Notifications []notifications.Notification
}

func listEvents(ctx context.Context, service di.Service, publicKey domain.PublicKey, tokensSet *internal.Set[domain.APNSToken]) error {
	fmt.Println()
	fmt.Println("listing stored events related to public key", publicKey.Hex())

	filters, err := domain.NewFilters(nostr.Filters{
		{
			Tags: map[string][]string{
				"p": {
					publicKey.Hex(),
				},
			},
			Limit: 100,
		},
	})
	if err != nil {
		return errors.Wrap(err, "error creating filters")
	}

	var events []eventWithNotifications
	var eventsLock sync.Mutex
	var eventsWaitGroup sync.WaitGroup

	for v := range service.App().Queries.GetEvents.Handle(ctx, filters) {
		if err := v.Err(); err != nil {
			return errors.Wrap(err, "handler returned an error")
		}

		if v.EOSE() {
			break
		}

		eventsWaitGroup.Add(1)

		event := v.Event()
		go func() {
			defer eventsWaitGroup.Done()

			notifications, err := service.App().Queries.GetNotifications.Handle(ctx, event.Id())
			if err != nil {
				fmt.Println(errors.Wrapf(err, "error getting notifications for event '%s'", event.Id().Hex()))
				return
			}

			eventsLock.Lock()
			defer eventsLock.Unlock()

			events = append(events, eventWithNotifications{
				Event:         event,
				Notifications: notifications,
			})
		}()
	}

	eventsWaitGroup.Wait()

	sort.Slice(events, func(i, j int) bool {
		return events[i].Event.CreatedAt().Before(events[j].Event.CreatedAt())
	})

	for _, eventWithNotifications := range events {
		evt := eventWithNotifications.Event

		fmt.Println("event", evt.Id().Hex(), "type", evt.Kind().Int(), "number of tags", len(evt.Tags()), "created at", evt.CreatedAt())

		if evt.PubKey() == publicKey {
			fmt.Println("-> own event")
		}

		for _, notification := range eventWithNotifications.Notifications {
			fmt.Printf("-> notification %s created at %s", notification.UUID(), notification.CreatedAt())
			if !tokensSet.Contains(notification.APNSToken()) {
				fmt.Print(" (for someone else)")
			}
			fmt.Println()
		}
	}

	return nil
}
