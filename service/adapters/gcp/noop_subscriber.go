package gcp

import (
	"context"

	"github.com/planetary-social/go-notification-service/service/domain"
)

type NoopSubscriber struct {
}

func NewNoopSubscriber() *NoopSubscriber {
	return &NoopSubscriber{}
}
func (p *NoopSubscriber) Subscribe(ctx context.Context) (<-chan *domain.FollowChange, error) {
	ch := make(chan *domain.FollowChange)
	return ch, nil
}
