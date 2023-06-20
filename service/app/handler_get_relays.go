package app

import (
	"context"

	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/domain"
)

type GetRelaysHandler struct {
	transactionProvider TransactionProvider
}

func NewGetRelaysHandler(
	transactionProvider TransactionProvider,
) *GetRelaysHandler {
	return &GetRelaysHandler{
		transactionProvider: transactionProvider,
	}
}

func (h *GetRelaysHandler) Handle(ctx context.Context) ([]domain.RelayAddress, error) {
	var result []domain.RelayAddress
	if err := h.transactionProvider.Transact(ctx, func(ctx context.Context, adapters Adapters) error {
		tmp, err := adapters.Registrations.GetRelays(ctx)
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
