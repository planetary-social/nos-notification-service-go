package app

import (
	"context"

	"github.com/planetary-social/go-notification-service/service/domain"
)

type SaveRegistration struct {
	eventAuthor  domain.PublicKey
	registration domain.Registration
}

func NewSaveRegistration(eventAuthor domain.PublicKey, registration domain.Registration) SaveRegistration {
	return SaveRegistration{eventAuthor: eventAuthor, registration: registration}
}

type SaveRegistrationHandler struct {
	transactionProvider TransactionProvider
}

func NewSaveRegistrationHandler(
	transactionProvider TransactionProvider,
) *SaveRegistrationHandler {
	return &SaveRegistrationHandler{
		transactionProvider: transactionProvider,
	}
}

func (h *SaveRegistrationHandler) Handle(ctx context.Context, cmd SaveRegistration) error {
	return h.transactionProvider.Transact(ctx, func(ctx context.Context, adapters Adapters) error {
		return adapters.Registrations.Save(cmd.registration)
	})
}
