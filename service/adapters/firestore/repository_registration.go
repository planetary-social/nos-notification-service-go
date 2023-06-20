package firestore

import (
	"encoding/hex"

	"cloud.google.com/go/firestore"
	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/domain"
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
)

func (r *RegistrationRepository) Save(registration domain.Registration) error {
	tokenDocPath := r.client.Collection(collectionAPNSTokens).Doc(registration.APNSToken().String())
	tokenDocData := map[string]interface{}{
		"token":  registration.APNSToken().String(),
		"locale": registration.Locale().String(),
	}
	if err := r.tx.Create(tokenDocPath, tokenDocData); err != nil {
		return errors.Wrap(err, "error creating the token doc")
	}

	for _, pubKeyWithRelays := range registration.PublicKeys() {
		pubKeyHex := hex.EncodeToString(pubKeyWithRelays.PublicKey().Bytes())

		publicKeyDocPath := tokenDocPath.Collection(collectionAPNSTokensPublicKeys).Doc(pubKeyHex)
		publicKeyDocData := map[string]interface{}{
			"publicKey": pubKeyHex,
		}
		if err := r.tx.Create(publicKeyDocPath, publicKeyDocData); err != nil {
			return errors.Wrap(err, "error creating the key doc")
		}

		for _, relayAddress := range pubKeyWithRelays.Relays() {
			relayDocPath := publicKeyDocPath.Collection(collectionAPNSTokensPublicKeysRelays).Doc(relayAddress.String())
			relayDocData := map[string]interface{}{
				"address": relayAddress.String(),
			}
			if err := r.tx.Create(relayDocPath, relayDocData); err != nil {
				return errors.Wrap(err, "error creating the key doc")
			}
		}
	}

	return nil
}
