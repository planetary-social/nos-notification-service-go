package firestore

import (
	"context"
	"encoding/hex"

	"cloud.google.com/go/firestore"
	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/domain"
	"google.golang.org/api/iterator"
)

const (
	collectionRelays           = "relays"
	collectionRelaysPublicKeys = "publicKeys"
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
		relayDocPath := r.client.Collection(collectionRelays).Doc(relayAddressAsKey(relayAddress))
		relayDocData := map[string]any{
			"address": relayAddress.String(),
		}
		if err := r.tx.Set(relayDocPath, relayDocData, firestore.MergeAll); err != nil {
			return errors.Wrap(err, "error creating the relay doc")
		}

		pubKeyDocPath := relayDocPath.Collection(collectionRelaysPublicKeys).Doc(registration.PublicKey().Hex())
		pubKeyDocData := map[string]any{
			"publicKey": registration.PublicKey().Hex(),
		}
		if err := r.tx.Set(pubKeyDocPath, pubKeyDocData, firestore.MergeAll); err != nil {
			return errors.Wrap(err, "error creating the public key doc")
		}
	}

	return nil
}

func (r *RelayRepository) GetRelays(ctx context.Context) ([]domain.RelayAddress, error) {
	// todo do it in transaction? emulator doesn't support it
	iter := r.client.Collection(collectionRelays).Documents(ctx)

	var result []domain.RelayAddress
	for {
		docRef, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, errors.Wrap(err, "error calling iter next")
		}

		relayAddress, err := relayAddressFromKey(docRef.Ref.ID)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating a relay address from key '%s'", docRef.Ref.ID)
		}
		result = append(result, relayAddress)
	}

	return result, nil
}

func (r *RelayRepository) GetPublicKeys(ctx context.Context, address domain.RelayAddress) ([]domain.PublicKey, error) {
	// todo do it in transaction? emulator doesn't support it
	iter := r.client.Collection(collectionRelays).Doc(relayAddressAsKey(address)).Collection(collectionRelaysPublicKeys).Documents(ctx)

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

func relayAddressAsKey(v domain.RelayAddress) string {
	return hex.EncodeToString([]byte(v.String()))
}

func relayAddressFromKey(v string) (domain.RelayAddress, error) {
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
