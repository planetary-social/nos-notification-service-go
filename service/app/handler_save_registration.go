package app

import (
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

func (h *SaveRegistrationHandler) Handle(cmd SaveRegistration) error {
	return h.transactionProvider.Transact(func(adapters Adapters) error {
		return adapters.Registrations.Save(cmd.registration)
	})
}
