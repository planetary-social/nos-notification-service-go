package mocks

import (
	"context"

	"github.com/planetary-social/go-notification-service/service/domain"
)

type MockExternalFollowChangeSubscriber struct {
}

func NewMockExternalFollowChangeSubscriber() *MockExternalFollowChangeSubscriber {
	return &MockExternalFollowChangeSubscriber{}
}

func (m MockExternalFollowChangeSubscriber) Subscribe(ctx context.Context) (<-chan *domain.FollowChangeBatch, error) {
	ch := make(chan *domain.FollowChangeBatch)
	return ch, nil
}
