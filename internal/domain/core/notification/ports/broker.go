package ports

import (
	"context"
	"wb-tech-l3/internal/domain/core/notification/model"
)

type Publisher interface {
	Publish(ctx context.Context, n *model.Notification) error
}

type Consumer interface {
	StartConsuming(
		ctx context.Context,
		handler func(ctx context.Context, m *model.Notification) error,
	) error
}

type MessagePipe interface {
	GetMessageChan() <-chan []byte
}

type Sender interface {
	PublishDelayed(notification model.Notification) error
}
