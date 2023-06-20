package app

import (
	"context"

	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/domain"
)

type GetPublicKeysHandler struct {
	transactionProvider TransactionProvider
}

func NewGetPublicKeysHandler(
	transactionProvider TransactionProvider,
) *GetPublicKeysHandler {
	return &GetPublicKeysHandler{
		transactionProvider: transactionProvider,
	}
}

func (h *GetPublicKeysHandler) Handle(ctx context.Context, relay domain.RelayAddress) ([]domain.PublicKey, error) {
	var result []domain.PublicKey
	if err := h.transactionProvider.Transact(ctx, func(ctx context.Context, adapters Adapters) error {
		tmp, err := adapters.Registrations.GetPublicKeys(ctx, relay)
		if err != nil {
			return errors.Wrap(err, "error getting relays")
		}
		result = tmp
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "transaction error")
	}
	return result, nil
}
