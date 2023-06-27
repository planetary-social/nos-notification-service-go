package app

import (
	"context"

	"github.com/planetary-social/go-notification-service/service/domain"
	"github.com/planetary-social/go-notification-service/service/domain/notifications"
)

type TransactionProvider interface {
	Transact(context.Context, func(context.Context, Adapters) error) error
}

type Adapters struct {
	Registrations RegistrationRepository
	Events        EventRepository
	Relays        RelayRepository
	PublicKeys    PublicKeyRepository
}

type RegistrationRepository interface {
	Save(registration domain.Registration) error
}

type RelayRepository interface {
	GetRelays(ctx context.Context) ([]domain.RelayAddress, error)
	GetPublicKeys(ctx context.Context, address domain.RelayAddress) ([]domain.PublicKey, error)
}

type PublicKeyRepository interface {
	GetAPNSTokens(context.Context, domain.PublicKey) ([]domain.APNSToken, error)
}

type EventRepository interface {
	Save(event domain.Event) error
	Exists(ctx context.Context, id domain.EventId) (bool, error)
	SaveNotificationForEvent(notification notifications.Notification) error
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
	GetTokens     *GetTokensHandler
}

type APNS interface {
	SendNotification(notification notifications.Notification) error
}
