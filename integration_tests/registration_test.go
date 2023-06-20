//go:build test_integration

package integration_tests

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/boreq/errors"
	"github.com/gorilla/websocket"
	"github.com/nbd-wtf/go-nostr"
	"github.com/planetary-social/go-notification-service/cmd/notification-service/di"
	"github.com/planetary-social/go-notification-service/internal/fixtures"
	"github.com/planetary-social/go-notification-service/service/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistration(t *testing.T) {
	ctx := fixtures.Context(t)
	config, service := createService(ctx, t)
	conn := createClient(ctx, t, config)

	publicKey, privateKeyHex := fixtures.SomeKeyPair()
	relayAddress := fixtures.SomeRelayAddress()

	event := nostr.Event{
		CreatedAt: nostr.Now(),
		Kind:      12345,
		Tags:      nostr.Tags{},
		Content: fmt.Sprintf(`
{
  "publicKeys": [
    {
      "publicKey": "%s",
      "relays": [
        {
          "address": "%s"
        }
      ]
    }
  ],
  "locale": "%s",
  "apnsToken": "%s"
}
`,
			publicKey.Hex(),
			relayAddress.String(),
			fixtures.SomeString(),
			fixtures.SomeString()),
	}

	err := event.Sign(privateKeyHex)
	require.NoError(t, err)

	envelope := nostr.EventEnvelope{
		SubscriptionID: nil,
		Event:          event,
	}

	j, err := envelope.MarshalJSON()
	require.NoError(t, err)

	err = conn.WriteMessage(websocket.TextMessage, j)
	require.NoError(t, err)

	require.EventuallyWithT(t, func(c *assert.CollectT) {
		relays, err := service.App().Queries.GetRelays.Handle(ctx)
		assert.NoError(c, err)
		assert.Contains(c, relays, relayAddress)
	}, 10*time.Second, 100*time.Millisecond)

	require.EventuallyWithT(t, func(c *assert.CollectT) {
		publicKeys, err := service.App().Queries.GetPublicKeys.Handle(ctx, relayAddress)
		assert.NoError(c, err)
		assert.Contains(c, publicKeys, publicKey)
	}, 10*time.Second, 100*time.Millisecond)
}

func createClient(ctx context.Context, tb testing.TB, config config.Config) *websocket.Conn {
	addr := config.NostrListenAddress()
	if strings.HasPrefix(addr, ":") {
		addr = "localhost" + addr
	}
	addr = fmt.Sprintf("ws://%s", addr)

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, addr, nil)
	require.NoError(tb, err)
	return conn
}

func createService(ctx context.Context, tb testing.TB) (config.Config, di.Service) {
	config, err := config.NewConfig("", "test-project-id")
	require.NoError(tb, err)

	service, cleanup, err := di.BuildService(ctx, config)
	require.NoError(tb, err)
	tb.Cleanup(cleanup)

	terminatedCh := make(chan error)

	tb.Cleanup(func() {
		if err := <-terminatedCh; err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			tb.Fatalf("error shutting down the service: %s", err)
		}
	})

	runCtx, cancelRunCtx := context.WithCancel(ctx)
	tb.Cleanup(cancelRunCtx)
	go func() {
		terminatedCh <- service.Run(runCtx)
	}()

	return config, service
}