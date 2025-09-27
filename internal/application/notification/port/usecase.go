package port

import (
	"context"
	"github.com/D1sordxr/url-shortener/internal/application/notification/input"
	"github.com/D1sordxr/url-shortener/internal/domain/core/notification/model"
)

type NotifyUseCase interface {
	Create(ctx context.Context, notify input.CreateNotifyInput) (*model.Notification, error)
	Read(ctx context.Context, id string) (*model.Notification, error)
	Delete(ctx context.Context, id string) error
}
