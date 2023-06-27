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
	GetEvents(ctx context.Context, filters domain.Filters) <-chan EventOrError
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
	GetEvents     *GetEventsHandler
}

type APNS interface {
	SendNotification(notification notifications.Notification) error
}

type EventOrError struct {
	event domain.Event
	err   error
}

func NewEventOrErrorWithEvent(event domain.Event) EventOrError {
	return EventOrError{event: event}
}

func NewEventOrErrorWithError(err error) EventOrError {
	return EventOrError{err: err}
}

func (e *EventOrError) Event() domain.Event {
	return e.event
}

func (e *EventOrError) Err() error {
	return e.err
}

type ReceivedEvent struct {
	relay domain.RelayAddress
	event domain.Event
}

func NewReceivedEvent(relay domain.RelayAddress, event domain.Event) ReceivedEvent {
	return ReceivedEvent{relay: relay, event: event}
}

func (r ReceivedEvent) Relay() domain.RelayAddress {
	return r.relay
}

func (r ReceivedEvent) Event() domain.Event {
	return r.event
}

type ReceivedEventSubscriber interface {
	Subscribe(ctx context.Context) <-chan ReceivedEvent
}
