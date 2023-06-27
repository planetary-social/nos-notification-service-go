package app

import (
	"context"

	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/service/domain"
)

type GetTokensHandler struct {
	transactionProvider TransactionProvider
}

func NewGetTokensHandler(
	transactionProvider TransactionProvider,
) *GetTokensHandler {
	return &GetTokensHandler{
		transactionProvider: transactionProvider,
	}
}

func (h *GetTokensHandler) Handle(ctx context.Context, publicKey domain.PublicKey) ([]domain.APNSToken, error) {
	var result []domain.APNSToken
	if err := h.transactionProvider.Transact(ctx, func(ctx context.Context, adapters Adapters) error {
		tmp, err := adapters.PublicKeys.GetAPNSTokens(ctx, publicKey)
		if err != nil {
			return errors.Wrap(err, "error getting apns tokens")
		}
		result = tmp
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "transaction error")
	}
	return result, nil
}
