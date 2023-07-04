package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/domain"
	"google.golang.org/api/iterator"
)

const (
	collectionPublicKeys           = "publicKeys"
	collectionPublicKeysAPNSTokens = "apnsTokens"
)

type PublicKeyRepository struct {
	client *firestore.Client
	tx     *firestore.Transaction
}

func NewPublicKeyRepository(client *firestore.Client, tx *firestore.Transaction) *PublicKeyRepository {
	return &PublicKeyRepository{client: client, tx: tx}
}

func (r *PublicKeyRepository) Save(registration domain.Registration) error {
	pubKeyDocPath := r.client.Collection(collectionPublicKeys).Doc(registration.PublicKey().Hex())
	pubKeyDocData := map[string]any{
		"publicKey": ensureType[string](registration.PublicKey().Hex()),
	}
	if err := r.tx.Set(pubKeyDocPath, pubKeyDocData, firestore.MergeAll); err != nil {
		return errors.Wrap(err, "error creating the public key doc")
	}

	tokenDocPath := r.client.Collection(collectionPublicKeys).Doc(registration.PublicKey().Hex()).Collection(collectionPublicKeysAPNSTokens).Doc(registration.APNSToken().Hex())
	tokenDocData := map[string]any{
		"token": ensureType[string](registration.APNSToken().Hex()),
	}
	if err := r.tx.Set(tokenDocPath, tokenDocData, firestore.MergeAll); err != nil {
		return errors.Wrap(err, "error creating the public key doc")
	}

	return nil
}

func (r *PublicKeyRepository) GetAPNSTokens(ctx context.Context, key domain.PublicKey) ([]domain.APNSToken, error) {
	// todo do it in transaction? emulator doesn't support it
	docs := r.client.Collection(collectionPublicKeys).Doc(key.Hex()).Collection(collectionPublicKeysAPNSTokens).Documents(ctx)

	var result []domain.APNSToken

	for {
		doc, err := docs.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, errors.Wrap(err, "error getting a document")
		}

		data := make(map[string]any)

		if err := doc.DataTo(&data); err != nil {
			return nil, errors.Wrap(err, "error reading document data")
		}

		apnsToken, err := domain.NewAPNSTokenFromHex(data["token"].(string))
		if err != nil {
			return nil, errors.Wrap(err, "error creating a token from hex")
		}

		result = append(result, apnsToken)
	}

	return result, nil
}
