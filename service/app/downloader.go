package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/boreq/errors"
	"github.com/gorilla/websocket"
	"github.com/nbd-wtf/go-nostr"
	"github.com/planetary-social/go-notification-service/internal"
	"github.com/planetary-social/go-notification-service/service/domain"
)

type Downloader struct {
	transactionProvider TransactionProvider
	relayDownloaders    map[domain.RelayAddress]*RelayDownloader
}

func NewDownloader(transaction TransactionProvider) *Downloader {
	return &Downloader{
		transactionProvider: transaction,
		relayDownloaders:    map[domain.RelayAddress]*RelayDownloader{},
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
				delete(d.relayDownloaders, relayAddress)
				relayDownloader.Stop()
			}
		}

		for _, relayAddress := range relayAddresses.List() {
			if _, ok := d.relayDownloaders[relayAddress]; !ok {
				relayDownloader := NewRelayDownloader(ctx, d.transactionProvider, relayAddress)
				d.relayDownloaders[relayAddress] = relayDownloader
			}
		}

		<-time.After(60 * time.Second)
	}
}

func (d *Downloader) getRelays(ctx context.Context) (*internal.Set[domain.RelayAddress], error) {
	var relays []domain.RelayAddress

	if err := d.transactionProvider.Transact(ctx, func(ctx context.Context, adapters Adapters) error {
		tmp, err := adapters.Registrations.GetRelays(ctx)
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
	address             domain.RelayAddress
	transactionProvider TransactionProvider
	cancel              context.CancelFunc
}

func NewRelayDownloader(ctx context.Context, transactionProvider TransactionProvider, address domain.RelayAddress) *RelayDownloader {
	ctx, cancel := context.WithCancel(ctx)
	v := &RelayDownloader{
		transactionProvider: transactionProvider,
		cancel:              cancel,
		address:             address,
	}
	go v.run(ctx)
	return v
}

func (d *RelayDownloader) run(ctx context.Context) {
	for {
		if err := d.connectAndDownload(ctx); err != nil {
			fmt.Printf("error processing relay '%s': %s\n", d.address, err)
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(10 * time.Second):
			continue
		}
	}
}

func (d *RelayDownloader) connectAndDownload(ctx context.Context) error {
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, d.address.String(), nil)
	if err != nil {
		return errors.Wrap(err, "error dialing the relay")
	}
	defer conn.Close()

	activeSubscriptions := internal.NewEmptySet[domain.PublicKey]()
	activeSubscriptionsLock := &sync.Mutex{}

	go func() {
		if err := d.manageSubs(ctx, conn, activeSubscriptions, activeSubscriptionsLock); err != nil {
			fmt.Println("error managing subs", err)
		}
	}()

	for {
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			return errors.Wrap(err, "error reading a message")
		}

		if err := d.handleMessage(ctx, messageBytes, activeSubscriptions, activeSubscriptionsLock); err != nil {
			return errors.Wrap(err, "error handling message")
		}

	}

}

func (d *RelayDownloader) handleMessage(
	ctx context.Context,
	messageBytes []byte,
	activeSubscriptions *internal.Set[domain.PublicKey],
	activeSubscriptionsLock *sync.Mutex,
) error {
	envelope := nostr.ParseMessage(messageBytes)
	if envelope == nil {
		return errors.New("error parsing message, we are never going to find out what error unfortunately due to the design of this library")
	}

	switch v := envelope.(type) {
	case *nostr.EOSEEnvelope:
		publicKey, err := domain.NewPublicKey(string(*v))
		if err != nil {
			return errors.Wrap(err, "invalid public key; unexpected subscription id since we only create them from public keys")
		}

		activeSubscriptionsLock.Lock()
		activeSubscriptionsLock.Unlock()
		activeSubscriptions.Delete(publicKey)
		// todo there is a bug here, we may have recreated the sub and this
		// message refers to the previous sub
	case *nostr.EventEnvelope:
		event, err := domain.NewEvent(v.Event)
		if err != nil {
			return errors.Wrap(err, "error creating an event")
		}

		// todo maybe pubsub those events and then handle them later?
		if err := d.transactionProvider.Transact(ctx, func(ctx context.Context, adapters Adapters) error {
			// todo figure out if we actually want to save this
			return adapters.Events.Save(d.address, event)
		}); err != nil {
			return errors.Wrap(err, "transaction error")
		}
	default:
		fmt.Println("unknown message:", string(messageBytes))
	}

	return nil
}

func (d *RelayDownloader) manageSubs(
	ctx context.Context,
	conn *websocket.Conn,
	activeSubscriptions *internal.Set[domain.PublicKey],
	activeSubscriptionsLock *sync.Mutex,
) error {
	defer conn.Close()

	for {
		publicKeys, err := d.getPublicKeys(ctx)
		if err != nil {
			return errors.Wrap(err, "error getting public keys")
		}

		if err := d.updateSubs(conn, activeSubscriptions, activeSubscriptionsLock, publicKeys); err != nil {
			return errors.Wrap(err, "error updating subscriptions")
		}
	}
}

func (d *RelayDownloader) updateSubs(
	conn *websocket.Conn,
	activeSubscriptions *internal.Set[domain.PublicKey],
	activeSubscriptionsLock *sync.Mutex,
	publicKeys *internal.Set[domain.PublicKey],
) error {
	activeSubscriptionsLock.Lock()
	defer activeSubscriptionsLock.Unlock()

	for _, publicKey := range activeSubscriptions.List() {
		if !publicKeys.Contains(publicKey) {
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
			envelope := nostr.ReqEnvelope{
				SubscriptionID: publicKey.Hex(),
				Filters: nostr.Filters{nostr.Filter{
					Authors: []string{
						publicKey.Hex(),
					},
					Since: nil, // todo filter based on already received events
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
		tmp, err := adapters.Registrations.GetPublicKeys(ctx, d.address)
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

func (d RelayDownloader) Stop() {
	d.cancel()
}
