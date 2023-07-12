package app

import (
	"context"

	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/domain"
)

type GetRelaysHandler struct {
	transactionProvider TransactionProvider
	metrics             Metrics
}

func NewGetRelaysHandler(
	transactionProvider TransactionProvider,
	metrics Metrics,
) *GetRelaysHandler {
	return &GetRelaysHandler{
		transactionProvider: transactionProvider,
		metrics:             metrics,
	}
}

func (h *GetRelaysHandler) Handle(ctx context.Context) (addresses []domain.RelayAddress, err error) {
	defer h.metrics.StartApplicationCall("getRelays").End(&err)

	var result []domain.RelayAddress
	if err := h.transactionProvider.Transact(ctx, func(ctx context.Context, adapters Adapters) error {
		tmp, err := adapters.Relays.GetRelays(ctx)
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
