package app

import (
	"context"
	"fmt"
	"time"

	"github.com/boreq/errors"
	"github.com/gorilla/websocket"
	"github.com/nbd-wtf/go-nostr"
	"github.com/planetary-social/go-notification-service/internal"
	"github.com/planetary-social/go-notification-service/internal/logging"
	"github.com/planetary-social/go-notification-service/service/domain"
)

const (
	reconnectEvery           = 1 * time.Minute
	manageSubscriptionsEvery = 1 * time.Minute

	howFarIntoThePastToLook = 365 * 24 * time.Hour
)

type ReceivedEventPublisher interface {
	Publish(relay domain.RelayAddress, event domain.Event)
}

type Downloader struct {
	transactionProvider    TransactionProvider
	receivedEventPublisher ReceivedEventPublisher
	logger                 logging.Logger

	relayDownloaders map[domain.RelayAddress]*RelayDownloader
}

func NewDownloader(
	transaction TransactionProvider,
	receivedEventPublisher ReceivedEventPublisher,
	logger logging.Logger,
) *Downloader {
	return &Downloader{
		transactionProvider:    transaction,
		receivedEventPublisher: receivedEventPublisher,
		logger:                 logger.New("downloader"),

		relayDownloaders: map[domain.RelayAddress]*RelayDownloader{},
	}
}

func (d *Downloader) Run(ctx context.Context) error {
	for {
		relayAddresses, err := d.getRelays(ctx)
		if err != nil {
			return errors.Wrap(err, "error getting relays")
		}

		for relayAddress, relayDownloader := range d.relayDownloaders {
			if !relayAddresses.Contains(relayAddress) {
				d.logger.Debug().
					WithField("relay", relayAddress.String()).
					Message("deleting a relay downloader")
				delete(d.relayDownloaders, relayAddress)
				relayDownloader.Stop()
			}
		}

		for _, relayAddress := range relayAddresses.List() {
			if _, ok := d.relayDownloaders[relayAddress]; !ok {
				d.logger.Debug().
					WithField("relay", relayAddress.String()).
					Message("creating a relay downloader")
				relayDownloader := NewRelayDownloader(
					ctx,
					d.transactionProvider,
					d.receivedEventPublisher,
					d.logger,
					relayAddress,
				)
				d.relayDownloaders[relayAddress] = relayDownloader
			}
		}

		<-time.After(60 * time.Second)
	}
}

func (d *Downloader) getRelays(ctx context.Context) (*internal.Set[domain.RelayAddress], error) {
	var relays []domain.RelayAddress

	if err := d.transactionProvider.Transact(ctx, func(ctx context.Context, adapters Adapters) error {
		tmp, err := adapters.Relays.GetRelays(ctx)
		if err != nil {
			return errors.Wrap(err, "error getting relays")
		}
		relays = tmp
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "transaction error")
	}

	return internal.NewSet(relays), nil
}

type RelayDownloader struct {
	transactionProvider    TransactionProvider
	receivedEventPublisher ReceivedEventPublisher
	logger                 logging.Logger

	address domain.RelayAddress
	cancel  context.CancelFunc
}

func NewRelayDownloader(
	ctx context.Context,
	transactionProvider TransactionProvider,
	receivedEventPublisher ReceivedEventPublisher,
	logger logging.Logger,
	address domain.RelayAddress,
) *RelayDownloader {
	ctx, cancel := context.WithCancel(ctx)
	v := &RelayDownloader{
		transactionProvider:    transactionProvider,
		receivedEventPublisher: receivedEventPublisher,
		logger:                 logger.New(fmt.Sprintf("relayDownloader(%s)", address)),

		cancel:  cancel,
		address: address,
	}
	go v.run(ctx)
	return v
}

func (d *RelayDownloader) run(ctx context.Context) {
	for {
		if err := d.connectAndDownload(ctx); err != nil {
			d.logger.Error().
				WithError(err).
				Message("error connecting and downloading")
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(reconnectEvery):
			continue
		}
	}
}

func (d *RelayDownloader) connectAndDownload(ctx context.Context) error {
	d.logger.Debug().Message("connecting")

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, d.address.String(), nil)
	if err != nil {
		return errors.Wrap(err, "error dialing the relay")
	}
	defer conn.Close()

	go func() {
		if err := d.manageSubs(ctx, conn); err != nil {
			d.logger.Error().
				WithError(err).
				Message("error managing subs")
		}
	}()

	for {
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			return errors.Wrap(err, "error reading a message")
		}

		if err := d.handleMessage(messageBytes); err != nil {
			return errors.Wrap(err, "error handling message")
		}

	}

}

func (d *RelayDownloader) handleMessage(messageBytes []byte) error {
	envelope := nostr.ParseMessage(messageBytes)
	if envelope == nil {
		return errors.New("error parsing message, we are never going to find out what error unfortunately due to the design of this library")
	}

	switch v := envelope.(type) {
	case *nostr.EOSEEnvelope:
		d.logger.Debug().
			WithField("subscription", string(*v)).
			Message("received EOSE")
	case *nostr.EventEnvelope:
		event, err := domain.NewEvent(v.Event)
		if err != nil {
			return errors.Wrap(err, "error creating an event")
		}

		d.receivedEventPublisher.Publish(d.address, event)
	default:
		d.logger.
			Error().
			WithField("message", string(messageBytes)).
			Message("unhandled message")
	}

	return nil
}

func (d *RelayDownloader) manageSubs(
	ctx context.Context,
	conn *websocket.Conn,
) error {
	defer conn.Close()

	activeSubscriptions := internal.NewEmptySet[domain.PublicKey]()

	for {
		publicKeys, err := d.getPublicKeys(ctx)
		if err != nil {
			return errors.Wrap(err, "error getting public keys")
		}

		if err := d.updateSubs(conn, activeSubscriptions, publicKeys); err != nil {
			return errors.Wrap(err, "error updating subscriptions")
		}

		select {
		case <-time.After(manageSubscriptionsEvery):
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (d *RelayDownloader) updateSubs(
	conn *websocket.Conn,
	activeSubscriptions *internal.Set[domain.PublicKey],
	publicKeys *internal.Set[domain.PublicKey],
) error {
	for _, publicKey := range activeSubscriptions.List() {
		if !publicKeys.Contains(publicKey) {
			d.logger.Debug().
				WithField("publicKey", publicKey).
				Message("closing subscription")

			envelope := nostr.CloseEnvelope(publicKey.Hex())

			envelopeJSON, err := envelope.MarshalJSON()
			if err != nil {
				return errors.Wrap(err, "marshaling close envelope failed")
			}

			if err := conn.WriteMessage(websocket.TextMessage, envelopeJSON); err != nil {
				return errors.Wrap(err, "writing close envelope error")
			}

			activeSubscriptions.Delete(publicKey)
		}
	}

	for _, publicKey := range publicKeys.List() {
		if ok := activeSubscriptions.Contains(publicKey); !ok {
			d.logger.Debug().
				WithField("publicKey", publicKey).
				Message("opening subscription")

			t := nostr.Timestamp(time.Now().Add(-howFarIntoThePastToLook).Unix())

			envelope := nostr.ReqEnvelope{
				SubscriptionID: publicKey.Hex(),
				Filters: nostr.Filters{nostr.Filter{
					Tags: map[string][]string{
						"p": {publicKey.Hex()},
					},
					Since: &t,
				}},
			}

			envelopeJSON, err := envelope.MarshalJSON()
			if err != nil {
				return errors.Wrap(err, "marshaling req envelope failed")
			}

			if err := conn.WriteMessage(websocket.TextMessage, envelopeJSON); err != nil {
				return errors.Wrap(err, "writing req envelope error")
			}

			activeSubscriptions.Put(publicKey)
		}
	}

	return nil
}

func (d *RelayDownloader) getPublicKeys(ctx context.Context) (*internal.Set[domain.PublicKey], error) {
	var publicKeys []domain.PublicKey

	if err := d.transactionProvider.Transact(ctx, func(ctx context.Context, adapters Adapters) error {
		tmp, err := adapters.Relays.GetPublicKeys(ctx, d.address)
		if err != nil {
			return errors.Wrap(err, "error getting public keys")
		}
		publicKeys = tmp
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "transaction error")
	}

	return internal.NewSet(publicKeys), nil
}

func (d *RelayDownloader) Stop() {
	d.cancel()
}
