package port

import (
	"context"
	"github.com/D1sordxr/url-shortener/internal/application/url/input"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/model"
)

type UseCase interface {
	Create(ctx context.Context, create input.Create) (*model.URL, error)
}
