package app

import (
	"context"

	"github.com/boreq/errors"
	"github.com/planetary-social/go-notification-service/internal/logging"
	"github.com/planetary-social/go-notification-service/service/domain"
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
	logger              logging.Logger
}

func NewProcessReceivedEventHandler(
	transactionProvider TransactionProvider,
	logger logging.Logger,
) *ProcessReceivedEventHandler {
	return &ProcessReceivedEventHandler{
		transactionProvider: transactionProvider,
		logger:              logger.New("processReceivedEventHandler"),
	}
}

func (h *ProcessReceivedEventHandler) Handle(ctx context.Context, cmd ProcessReceivedEvent) error {
	h.logger.Debug().
		WithField("relay", cmd.relay.String()).
		WithField("event.id", cmd.event.Id().Hex()).
		WithField("event.kind", cmd.event.Kind().Int()).
		Message("processing received event")

	if err := h.transactionProvider.Transact(ctx, func(ctx context.Context, adapters Adapters) error {
		// todo figure out if we actually want to save this?

		return adapters.Events.Save(cmd.relay, cmd.event)
	}); err != nil {
		return errors.Wrap(err, "transaction error")
	}

	return nil
}
