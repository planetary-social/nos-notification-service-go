package firestore

import (
	"cloud.google.com/go/firestore"
	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/domain"
)

type RegistrationRepository struct {
	client *firestore.Client
	tx     *firestore.Transaction

	relayRepository     *RelayRepository
	publicKeyRepository *PublicKeyRepository
}

func NewRegistrationRepository(
	client *firestore.Client,
	tx *firestore.Transaction,
	relayRepository *RelayRepository,
	publicKeyRepository *PublicKeyRepository,
) *RegistrationRepository {
	return &RegistrationRepository{
		client:              client,
		tx:                  tx,
		relayRepository:     relayRepository,
		publicKeyRepository: publicKeyRepository,
	}
}

const (
	collectionAPNSTokens                 = "apnsTokens"
	collectionAPNSTokensPublicKeys       = "publicKeys"
	collectionAPNSTokensPublicKeysRelays = "relays"
)

func (r *RegistrationRepository) Save(registration domain.Registration) error {
	if err := r.saveUnderTokens(registration); err != nil {
		return errors.Wrap(err, "error saving under tokens")
	}

	if err := r.relayRepository.Save(registration); err != nil {
		return errors.Wrap(err, "error saving under relays")
	}

	if err := r.publicKeyRepository.Save(registration); err != nil {
		return errors.Wrap(err, "error saving under relays")
	}

	return nil
}

func (r *RegistrationRepository) saveUnderTokens(registration domain.Registration) error {
	tokenDocPath := r.client.Collection(collectionAPNSTokens).Doc(registration.APNSToken().Hex())
	tokenDocData := map[string]any{
		"token": registration.APNSToken().Hex(),
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
			relayDocPath := publicKeyDocPath.Collection(collectionAPNSTokensPublicKeysRelays).Doc(relayAddressAsKey(relayAddress))
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
