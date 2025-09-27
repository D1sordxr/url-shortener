package ports

import (
	"context"
	"wb-tech-l3/internal/domain/core/notification/model"
)

type CacheStore interface {
	Create(ctx context.Context, notification *model.Notification) error
	Read(ctx context.Context, id string) (*model.Notification, error)
	Delete(ctx context.Context, id string) error
	SetDeleted(ctx context.Context, id string) error
	IsDeleted(ctx context.Context, id string) (bool, error)
}
