package ports

import (
	"context"

	"github.com/D1sordxr/url-shortener/internal/domain/core/url/model"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/params"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/vo"
)

type Repository interface {
	Create(ctx context.Context, p params.CreateURL) (*model.URL, error)
	GetCompleteAnalytics(ctx context.Context, alias vo.Alias) (*model.Analytics, error)
	ReadURLByAlias(ctx context.Context, alias vo.Alias) (*model.URL, error)
	CreateStat(ctx context.Context, p params.CreateURLStat) (*model.URLStat, error)
}
