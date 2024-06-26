package firestore

import (
	"context"
	"encoding/hex"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/domain"
	"google.golang.org/api/iterator"
)

var relaySuffixesToSkip = []string{
	"127.0.0.1",
	"localhost",
	"nostr.band",
	"nostrja-kari-nip50.heguro.com",
	"nostr.sebastix.social",
}

const (
	collectionRelays                      = "relays"
	collectionRelaysFieldAddress          = "address"
	collectionRelaysFieldUpdatedTimestamp = "updatedTimestamp"

	collectionRelaysPublicKeys                      = "publicKeys"
	collectionRelaysPublicKeysFieldPublicKey        = "publicKey"
	collectionRelaysPublicKeysFieldUpdatedTimestamp = "updatedTimestamp"
)

type RelayRepository struct {
	client *firestore.Client
	tx     *firestore.Transaction
}

func NewRelayRepository(client *firestore.Client, tx *firestore.Transaction) *RelayRepository {
	return &RelayRepository{client: client, tx: tx}
}

func (r *RelayRepository) Save(registration domain.Registration) error {
	for _, relayAddress := range registration.Relays() {
		relayDocPath := r.client.Collection(collectionRelays).Doc(r.relayAddressAsKey(relayAddress))
		relayDocData := map[string]any{
			collectionRelaysFieldAddress:          ensureType[string](relayAddress.String()),
			collectionRelaysFieldUpdatedTimestamp: ensureType[time.Time](time.Now()),
		}
		if err := r.tx.Set(relayDocPath, relayDocData, firestore.MergeAll); err != nil {
			return errors.Wrap(err, "error creating the relay doc")
		}

		pubKeyDocPath := relayDocPath.Collection(collectionRelaysPublicKeys).Doc(registration.PublicKey().Hex())
		pubKeyDocData := map[string]any{
			collectionRelaysPublicKeysFieldPublicKey:        ensureType[string](registration.PublicKey().Hex()),
			collectionRelaysPublicKeysFieldUpdatedTimestamp: ensureType[time.Time](time.Now()),
		}
		if err := r.tx.Set(pubKeyDocPath, pubKeyDocData, firestore.MergeAll); err != nil {
			return errors.Wrap(err, "error creating the public key doc")
		}
	}

	return nil
}

func (r *RelayRepository) GetRelays(ctx context.Context, updatedAfter time.Time) ([]domain.RelayAddress, error) {
	iter := r.tx.Documents(
		r.client.
			Collection(collectionRelays).
			Where(collectionRelaysFieldUpdatedTimestamp, ">", updatedAfter),
	)

	var result []domain.RelayAddress
	for {
		docRef, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, errors.Wrap(err, "error calling iter next")
		}

		relayAddress, err := r.relayAddressFromKey(docRef.Ref.ID)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating a relay address from key '%s'", docRef.Ref.ID)
		}

		if !endsWithAny(relayAddress.HostWithoutPort(), relaySuffixesToSkip) {
			result = append(result, relayAddress)
		}
	}

	return result, nil
}

// endsWithAny checks if the given string ends with any of the strings in the
// list.
func endsWithAny(s string, list []string) bool {
	for _, suffix := range list {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}

func (r *RelayRepository) GetPublicKeys(ctx context.Context, address domain.RelayAddress, updatedAfter time.Time) ([]domain.PublicKey, error) {
	iter := r.tx.Documents(
		r.client.
			Collection(collectionRelays).
			Doc(r.relayAddressAsKey(address)).
			Collection(collectionRelaysPublicKeys).
			Where(collectionRelaysFieldUpdatedTimestamp, ">", updatedAfter),
	)

	var result []domain.PublicKey
	for {
		docRef, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, errors.Wrap(err, "error calling iter next")
		}

		publicKey, err := domain.NewPublicKeyFromHex(docRef.Ref.ID)
		if err != nil {
			return nil, errors.Wrap(err, "error creating a public key")
		}
		result = append(result, publicKey)
	}

	return result, nil
}

func (r *RelayRepository) relayAddressAsKey(v domain.RelayAddress) string {
	return hex.EncodeToString([]byte(v.String()))
}

func (r *RelayRepository) relayAddressFromKey(v string) (domain.RelayAddress, error) {
	b, err := hex.DecodeString(v)
	if err != nil {
		return domain.RelayAddress{}, errors.Wrap(err, "error decoding relay address from hex")
	}

	addr, err := domain.NewRelayAddress(string(b))
	if err != nil {
		return domain.RelayAddress{}, errors.Wrap(err, "error creating a relay address")
	}

	return addr, nil
}
