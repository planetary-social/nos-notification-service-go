package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/app"
	"github.com/planetary-social/go-notification-service/service/config"
	"google.golang.org/api/option"
)

func NewClient(ctx context.Context, config config.Config) (*firestore.Client, error) {
	var options []option.ClientOption

	if j := config.FirestoreCredentialsJSON(); len(j) > 0 {
		options = append(options, option.WithCredentialsJSON(config.FirestoreCredentialsJSON()))
	}

	return firestore.NewClient(ctx, config.FirestoreProjectID(), options...)
}

type AdaptersFactoryFn func(*firestore.Client, *firestore.Transaction) (app.Adapters, error)

type TransactionProvider struct {
	fn     AdaptersFactoryFn
	client *firestore.Client
}

func NewTransactionProvider(client *firestore.Client, fn AdaptersFactoryFn) *TransactionProvider {
	return &TransactionProvider{
		fn:     fn,
		client: client,
	}
}

func (t *TransactionProvider) Transact(ctx context.Context, f func(context.Context, app.Adapters) error) error {
	if err := t.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		adapters, err := t.fn(t.client, tx)
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
