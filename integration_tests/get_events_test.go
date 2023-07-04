//go:build test_integration

package integration_tests

import (
	"context"
	"testing"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/planetary-social/go-notification-service/internal/fixtures"
	"github.com/planetary-social/go-notification-service/service/app"
	"github.com/planetary-social/go-notification-service/service/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEvents(t *testing.T) {
	ctx := fixtures.Context(t)
	_, service := createService(ctx, t)

	_, sk1 := fixtures.SomeKeyPair()
	pk2, _ := fixtures.SomeKeyPair()

	libevent := nostr.Event{
		CreatedAt: nostr.Now(),
		Kind:      domain.EventKindNote.Int(),
		Tags: nostr.Tags{
			{"p", pk2.Hex()},
		},
		Content: "some content",
	}

	err := libevent.Sign(sk1)
	require.NoError(t, err)

	event, err := domain.NewEvent(libevent)
	require.NoError(t, err)

	cmd := app.NewSaveReceivedEvent(fixtures.SomeRelayAddress(), event)
	err = service.App().Commands.SaveReceivedEvent.Handle(ctx, cmd)
	require.NoError(t, err)

	require.EventuallyWithT(t, func(c *assert.CollectT) {
		since := nostr.Timestamp(time.Now().Add(-1 * time.Hour).Unix())

		filters, err := domain.NewFilters(nostr.Filters{
			{
				IDs:     nil,
				Kinds:   []int{domain.EventKindNote.Int(), domain.EventKindReaction.Int()},
				Authors: nil,
				Tags: nostr.TagMap{
					"p": []string{pk2.Hex()},
				},
				Since:  &since,
				Until:  nil,
				Limit:  0,
				Search: "",
			},
		})
		require.NoError(t, err)

		var events []domain.Event

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		for v := range service.App().Queries.GetEvents.Handle(ctx, filters) {
			if err := v.Err(); err != nil {
				c.Errorf("error: %s", err)
			}

			if v.EOSE() {
				break
			}

			events = append(events, v.Event())
		}

		assert.Len(t, events, 1)
		assert.Equal(t, event.Id(), events[0].Id())
	}, durationTimeout, durationTick)
}
