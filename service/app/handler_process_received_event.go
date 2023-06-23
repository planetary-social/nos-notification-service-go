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
		exists, err := adapters.Events.Exists(ctx, cmd.event.Id())
		if err != nil {
			return errors.Wrap(err, "error checking if event exists")
		}

		if exists {
			return nil
		}

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
				// todo send via pubsub instead
				if err := h.apns.SendNotification(notification); err != nil {
					return errors.Wrap(err, "error sending a notification")
				}

				if err := adapters.Events.SaveNotificationForEvent(notification); err != nil {
					return errors.Wrap(err, "error saving notification")
				}
			}
		}

		// todo don't save if we don't find this event relevant in the loop above?
		if err := adapters.Events.Save(cmd.event); err != nil {
			return errors.Wrap(err, "error saving the event")
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "transaction error")
	}

	return nil
}
