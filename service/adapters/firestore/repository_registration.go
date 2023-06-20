package firestore

import (
	"context"
	"encoding/hex"

	"cloud.google.com/go/firestore"
	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/domain"
	"google.golang.org/api/iterator"
)

type RegistrationRepository struct {
	client *firestore.Client
	tx     *firestore.Transaction
}

func NewRegistrationRepository(client *firestore.Client, tx *firestore.Transaction) *RegistrationRepository {
	return &RegistrationRepository{client: client, tx: tx}
}

const (
	collectionAPNSTokens                 = "apnsTokens"
	collectionAPNSTokensPublicKeys       = "publicKeys"
	collectionAPNSTokensPublicKeysRelays = "relays"

	collectionRelays           = "relays"
	collectionRelaysPublicKeys = "publicKeys"
)

func (r *RegistrationRepository) Save(registration domain.Registration) error {
	if err := r.saveUnderTokens(registration); err != nil {
		return errors.Wrap(err, "error saving under tokens")
	}

	if err := r.saveUnderRelays(registration); err != nil {
		return errors.Wrap(err, "error saving under relays")
	}

	return nil
}

func (r *RegistrationRepository) saveUnderTokens(registration domain.Registration) error {
	tokenDocPath := r.client.Collection(collectionAPNSTokens).Doc(registration.APNSToken().String())
	tokenDocData := map[string]any{
		"token":  registration.APNSToken().String(),
		"locale": registration.Locale().String(),
	}
	if err := r.tx.Set(tokenDocPath, tokenDocData, firestore.MergeAll); err != nil {
		return errors.Wrap(err, "error updating the token doc")
	}

	for _, pubKeyWithRelays := range registration.PublicKeys() {
		pubKeyHex := pubKeyWithRelays.PublicKey().Hex()

		publicKeyDocPath := tokenDocPath.Collection(collectionAPNSTokensPublicKeys).Doc(pubKeyHex)
		publicKeyDocData := map[string]any{
			"publicKey": pubKeyHex,
		}
		if err := r.tx.Set(publicKeyDocPath, publicKeyDocData, firestore.MergeAll); err != nil {
			return errors.Wrap(err, "error creating the key doc")
		}

		for _, relayAddress := range pubKeyWithRelays.Relays() {
			relayDocPath := publicKeyDocPath.Collection(collectionAPNSTokensPublicKeysRelays).Doc(r.relayAddressAsKey(relayAddress))
			relayDocData := map[string]any{
				"address": relayAddress.String(),
			}
			if err := r.tx.Set(relayDocPath, relayDocData, firestore.MergeAll); err != nil {
				return errors.Wrap(err, "error creating the key doc")
			}
		}
	}

	return nil
}

func (r *RegistrationRepository) saveUnderRelays(registration domain.Registration) error {
	for _, pubKeyWithRelays := range registration.PublicKeys() {
		for _, relayAddress := range pubKeyWithRelays.Relays() {
			relayDocPath := r.client.Collection(collectionRelays).Doc(r.relayAddressAsKey(relayAddress))
			relayDocData := map[string]any{
				"address": relayAddress,
			}
			if err := r.tx.Set(relayDocPath, relayDocData, firestore.MergeAll); err != nil {
				return errors.Wrap(err, "error creating the relay doc")
			}

			pubKeyDocPath := relayDocPath.Collection(collectionRelaysPublicKeys).Doc(pubKeyWithRelays.PublicKey().Hex())
			pubKeyDocData := map[string]any{
				"publicKey": pubKeyWithRelays.PublicKey().Hex(),
			}
			if err := r.tx.Set(pubKeyDocPath, pubKeyDocData, firestore.MergeAll); err != nil {
				return errors.Wrap(err, "error creating the public key doc")
			}
		}
	}
	return nil
}

func (r *RegistrationRepository) GetRelays(ctx context.Context) ([]domain.RelayAddress, error) {
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

		relayAddress, err := r.relayAddressFromKey(docRef.Ref.ID)
		if err != nil {
			return nil, errors.Wrapf(err, "error creating a relay address from key '%s'", docRef.Ref.ID)
		}
		result = append(result, relayAddress)
	}

	return result, nil
}

func (r *RegistrationRepository) GetPublicKeys(ctx context.Context, address domain.RelayAddress) ([]domain.PublicKey, error) {
	iter := r.client.Collection(collectionRelays).Doc(r.relayAddressAsKey(address)).Collection(collectionRelaysPublicKeys).Documents(ctx)

	var result []domain.PublicKey
	for {
		docRef, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, errors.Wrap(err, "error calling iter next")
		}

		publicKey, err := domain.NewPublicKey(docRef.Ref.ID)
		if err != nil {
			return nil, errors.Wrap(err, "error creating a public key")
		}
		result = append(result, publicKey)
	}

	return result, nil
}

func (r *RegistrationRepository) relayAddressAsKey(v domain.RelayAddress) string {
	return hex.EncodeToString([]byte(v.String()))
}

func (r *RegistrationRepository) relayAddressFromKey(v string) (domain.RelayAddress, error) {
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
