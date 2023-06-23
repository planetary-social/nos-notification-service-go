package app

import (
	"context"

	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/internal/logging"
	"github.com/planetary-social/go-notification-service/service/adapters/apns"
	"github.com/planetary-social/go-notification-service/service/domain"
	"github.com/planetary-social/go-notification-service/service/domain/notifications"
)

type ProcessReceivedEvent struct {
	relay domain.RelayAddress
	event domain.Event
}

func NewProcessReceivedEvent(relay domain.RelayAddress, event domain.Event) ProcessReceivedEvent {
	return ProcessReceivedEvent{relay: relay, event: event}
}

type ProcessReceivedEventHandler struct {
	transactionProvider TransactionProvider
	generator           *notifications.Generator
	apns                *apns.APNS
	logger              logging.Logger
}

func NewProcessReceivedEventHandler(
	transactionProvider TransactionProvider,
	generator *notifications.Generator,
	apns *apns.APNS,
	logger logging.Logger,
) *ProcessReceivedEventHandler {
	return &ProcessReceivedEventHandler{
		transactionProvider: transactionProvider,
		generator:           generator,
		apns:                apns,
		logger:              logger.New("processReceivedEventHandler"),
	}
}

func (h *ProcessReceivedEventHandler) Handle(ctx context.Context, cmd ProcessReceivedEvent) error {
	h.logger.Debug().
		WithField("relay", cmd.relay.String()).
		WithField("event.id", cmd.event.Id().Hex()).
		WithField("event.kind", cmd.event.Kind().Int()).
		Message("processing received event")

	mentions, err := domain.GetMentionsFromTags(cmd.event.Tags())
	if err != nil {
		return errors.Wrap(err, "error getting mentions for this event")
	}

	if err := h.transactionProvider.Transact(ctx, func(ctx context.Context, adapters Adapters) error {
		for _, mention := range mentions {
			token, err := adapters.PublicKeys.GetAPNSToken(ctx, mention)
			if err != nil {
				if errors.Is(err, APNSTokenNotFoundErr) {
					continue
				}
				return errors.Wrap(err, "error getting the token")
			}

			notifications, err := h.generator.Generate(mention, token, cmd.event)
			if err != nil {
				return errors.Wrap(err, "error generating notifications")
			}

			for _, notification := range notifications {
				id, err := h.apns.SendNotification(notification)
				if err != nil {
					return errors.Wrap(err, "error sending a notification")
				}

				h.logger.Debug().WithField("id", id.String()).Message("sent a notification")
			}
		}

		// todo maybe not always save?
		return adapters.Events.Save(cmd.relay, cmd.event)
	}); err != nil {
		return errors.Wrap(err, "transaction error")
	}

	return nil
}
