package app

import (
	"context"

	"github.com/planetary-social/go-notification-service/service/domain"
)

type TransactionProvider interface {
	Transact(context.Context, func(context.Context, Adapters) error) error
}

type Adapters struct {
	Registrations RegistrationRepository
	Events        EventRepository
	Relays        RelayRepository
}

type RegistrationRepository interface {
	Save(registration domain.Registration) error
}

type RelayRepository interface {
	GetRelays(ctx context.Context) ([]domain.RelayAddress, error)
	GetPublicKeys(ctx context.Context, address domain.RelayAddress) ([]domain.PublicKey, error)
}

type EventRepository interface {
	Save(relay domain.RelayAddress, event domain.Event) error
}

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	SaveRegistration *SaveRegistrationHandler
}

type Queries struct {
	GetRelays     *GetRelaysHandler
	GetPublicKeys *GetPublicKeysHandler
}
