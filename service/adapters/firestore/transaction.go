package firestore

import (
	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/app"
)

type TransactionProvider struct {
}

func NewTransactionProvider() *TransactionProvider {
	return &TransactionProvider{}
}

func (t TransactionProvider) Transact(f func(adapters app.Adapters) error) error {
	return errors.New("not implemented")
}
