package port

import (
	"context"
	"github.com/D1sordxr/url-shortener/internal/application/url/input"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/model"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/vo"
)

type UseCase interface {
	Create(ctx context.Context, create input.Create) (*model.URL, error)
	GetCompleteAnalytics(ctx context.Context, alias string) (*model.Analytics, error)
	GetForRedirect(ctx context.Context, i input.GetForRedirect) (vo.URL, error)
}
