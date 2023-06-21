package app

import (
	"context"

	"github.com/boreq/errors"
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
}

func NewProcessReceivedEventHandler(
	transactionProvider TransactionProvider,
) *ProcessReceivedEventHandler {
	return &ProcessReceivedEventHandler{
		transactionProvider: transactionProvider,
	}
}

func (h *ProcessReceivedEventHandler) Handle(ctx context.Context, cmd ProcessReceivedEvent) error {
	if err := h.transactionProvider.Transact(ctx, func(ctx context.Context, adapters Adapters) error {
		// todo figure out if we actually want to save this
		return adapters.Events.Save(cmd.relay, cmd.event)
	}); err != nil {
		return errors.Wrap(err, "transaction error")
	}

	return nil
}
