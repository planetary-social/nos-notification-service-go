package notifications_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/planetary-social/go-notification-service/internal/fixtures"
	"github.com/planetary-social/go-notification-service/internal/logging"
	"github.com/planetary-social/go-notification-service/service/domain"
	"github.com/planetary-social/go-notification-service/service/domain/notifications"
	"github.com/stretchr/testify/require"
)

func TestGenerator(t *testing.T) {
	testCases := []struct {
		Name string

		EventKind domain.EventKind

		ExpectedLocKey   string
		ExpectedCategory string
	}{
		{
			Name: "note",

			EventKind: domain.EventKindNote,

			ExpectedLocKey:   "NOTIFICATION_TAGGED_IN_NOTE",
			ExpectedCategory: "event.tagged.note",
		},
		{
			Name: "reaction",

			EventKind: domain.EventKindReaction,

			ExpectedLocKey:   "NOTIFICATION_TAGGED_IN_REACTION",
			ExpectedCategory: "event.tagged.reaction",
		},
		{
			Name: "edm",

			EventKind: domain.EventKindEncryptedDirectMessage,

			ExpectedLocKey:   "NOTIFICATION_TAGGED_IN_ENCRYPTED_DIRECT_MESSAGE",
			ExpectedCategory: "event.tagged.encryptedDirectMessage",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			logger := logging.NewDevNullLogger()
			g := notifications.NewGenerator(logger)

			pk1, _ := fixtures.SomeKeyPair()
			pk2, sk2 := fixtures.SomeKeyPair()

			libevent := nostr.Event{
				PubKey:    pk2.Hex(),
				CreatedAt: nostr.Timestamp(time.Now().Unix()),
				Kind:      testCase.EventKind.Int(),
				Tags: nostr.Tags{
					nostr.Tag{"p", pk1.Hex()},
				},
				Content: "some content",
			}

			err := libevent.Sign(sk2)
			require.NoError(t, err)

			event, err := domain.NewEvent(libevent)
			require.NoError(t, err)

			token := fixtures.SomeAPNSToken()

			result, err := g.Generate(pk1, token, event)
			require.NoError(t, err)

			require.Len(t, result, 1)

			notification := result[0]
			require.Equal(t,
				fmt.Sprintf(
					`{"aps":{"alert":{"loc-key":"%s"},"category":"%s"},"event":{"id":"%s","pubkey":"%s","created_at":%d,"kind":%d,"tags":[["p","%s"]],"content":"some content","sig":"%s"}}`,
					testCase.ExpectedLocKey,
					testCase.ExpectedCategory,
					event.Id().Hex(),
					pk2.Hex(),
					event.CreatedAt().Unix(),
					testCase.EventKind.Int(),
					pk1.Hex(),
					event.Sig().Hex(),
				),
				string(notification.Payload()),
			)
			require.Equal(t, token, notification.APNSToken())
			require.Equal(t, event, notification.Event())
		})
	}
}
