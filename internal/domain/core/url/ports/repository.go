package ports

import (
	"context"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/model"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/params"
)

type Repository interface {
	Create(ctx context.Context, p params.CreateURL) (*model.URL, error)
	GetAnalyticsByAlias(ctx context.Context, alias string) ([]model.URLStat, error)
}
