package app

import "github.com/planetary-social/go-notification-service/service/domain"

type TransactionProvider interface {
	Transact(func(adapters Adapters) error) error
}

type Adapters struct {
	Registrations RegistrationRepository
}

type RegistrationRepository interface {
	Save(registration domain.Registration) error
}
