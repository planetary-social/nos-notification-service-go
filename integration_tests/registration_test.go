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
	"github.com/stretchr/testify/require"
)

func TestRegistration(t *testing.T) {
	ctx := fixtures.Context(t)
	config, _ := createService(ctx, t)
	conn := createClient(ctx, t, config)

	privateKey := nostr.GeneratePrivateKey()

	event := nostr.Event{
		CreatedAt: nostr.Now(),
		Kind:      12345,
		Tags:      nostr.Tags{},
		Content: `
{
  "publicKeys": [
    {
      "publicKey": "some-public-key",
      "relays": [
        {
          "address": "some-relay-address"
        }
      ]
    }
  ],
  "locale": "some-locale",
  "apnsToken": "some-apns-token"
}
`,
	}
	err := event.Sign(privateKey)
	require.NoError(t, err)

	envelope := nostr.EventEnvelope{
		SubscriptionID: nil,
		Event:          event,
	}

	j, err := envelope.MarshalJSON()
	require.NoError(t, err)

	err = conn.WriteMessage(websocket.TextMessage, j)
	require.NoError(t, err)

	<-time.After(1 * time.Second) // todo replace with some kind a success condition
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
	config := config.NewConfig("")
	service, cleanup, err := di.BuildService(config)
	require.NoError(tb, err)
	tb.Cleanup(cleanup)

	terminatedCh := make(chan error)

	runCtx, cancelRunCtx := context.WithCancel(ctx)

	tb.Cleanup(func() {
		if err := <-terminatedCh; err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			tb.Fatalf("error shutting down the service: %s", err)
		}
	})

	tb.Cleanup(cancelRunCtx)

	go func() {
		terminatedCh <- service.Run(runCtx)
	}()

	return config, service
}
