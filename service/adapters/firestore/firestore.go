package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/app"
	"github.com/planetary-social/go-notification-service/service/config"
)

func NewClient(ctx context.Context, config config.Config) (*firestore.Client, error) {
	return firestore.NewClient(ctx, config.FirestoreProjectID())
}

type AdaptersFactoryFn func(*firestore.Transaction) (app.Adapters, error)

type TransactionProvider struct {
	fn     AdaptersFactoryFn
	client *firestore.Client
}

func NewTransactionProvider(client *firestore.Client) *TransactionProvider {
	return &TransactionProvider{
		client: client,
	}
}

func (t *TransactionProvider) Transact(ctx context.Context, f func(context.Context, app.Adapters) error) error {
	if err := t.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		adapters, err := t.fn(tx)
		if err != nil {
			return errors.Wrap(err, "error building the adapters")
		}

		if err := f(ctx, adapters); err != nil {
			return errors.Wrap(err, "error calling the provided function")
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "transaction returned an error")
	}

	return nil
}
