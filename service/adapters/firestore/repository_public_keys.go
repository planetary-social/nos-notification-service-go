package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/app"
	"github.com/planetary-social/go-notification-service/service/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	collectionPublicKeys = "publicKeys"
)

type PublicKeyRepository struct {
	client *firestore.Client
	tx     *firestore.Transaction
}

func NewPublicKeyRepository(client *firestore.Client, tx *firestore.Transaction) *PublicKeyRepository {
	return &PublicKeyRepository{client: client, tx: tx}
}

func (r *PublicKeyRepository) Save(registration domain.Registration) error {
	for _, pubKeyWithRelays := range registration.PublicKeys() {
		pubKeyDocPath := r.client.Collection(collectionPublicKeys).Doc(pubKeyWithRelays.PublicKey().Hex())
		pubKeyDocData := map[string]any{
			"publicKey": pubKeyWithRelays.PublicKey().Hex(),
			"token":     registration.APNSToken().Hex(),
		}
		if err := r.tx.Set(pubKeyDocPath, pubKeyDocData, firestore.MergeAll); err != nil {
			return errors.Wrap(err, "error creating the public key doc")
		}
	}
	return nil
}

func (r *PublicKeyRepository) GetAPNSToken(ctx context.Context, key domain.PublicKey) (domain.APNSToken, error) {
	// todo do it in transaction? emulator doesn't support it
	doc, err := r.client.Collection(collectionPublicKeys).Doc(key.Hex()).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.APNSToken{}, app.APNSTokenNotFoundErr
		}
		return domain.APNSToken{}, errors.Wrap(err, "error getting the public key document")
	}

	data := make(map[string]any)

	if err := doc.DataTo(&data); err != nil {
		return domain.APNSToken{}, errors.Wrap(err, "error reading document data")
	}

	apnsToken, err := domain.NewAPNSTokenFromHex(data["token"].(string))
	if err != nil {
		return domain.APNSToken{}, errors.Wrap(err, "error creating a token from hex")
	}

	return apnsToken, nil
}
