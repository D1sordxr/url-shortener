package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	appPorts "github.com/D1sordxr/url-shortener/internal/domain/app/ports"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/errorx"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/model"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/params"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/vo"
	"github.com/D1sordxr/url-shortener/internal/infrastructure/storage/postgres/errordb"
	"github.com/D1sordxr/url-shortener/internal/infrastructure/storage/postgres/repositories/url/converters"
	"github.com/D1sordxr/url-shortener/internal/infrastructure/storage/postgres/repositories/url/gen"
	"github.com/wb-go/wbf/dbpg"
)

type Repository struct {
	log      appPorts.Logger
	executor *dbpg.DB
	queries  *gen.Queries
}

func NewRepository(log appPorts.Logger, executor *dbpg.DB) *Repository {
	return &Repository{
		log:      log,
		executor: executor,
		queries:  gen.New(executor.Master),
	}
}

func (r *Repository) Create(ctx context.Context, p params.CreateURL) (*model.URL, error) {
	const op = "postgres.url.Repository.Create"

	rawUrl, err := r.queries.CreateURL(ctx, converters.ConvertCreateParams(p))
	if err != nil {
		if errordb.IsUniqueViolation(err) {
			return nil, fmt.Errorf("%s: %w", op, errorx.ErrAliasAlreadyExists)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	url := converters.ConvertGenToDomain(rawUrl)
	return &url, nil
}

func (r *Repository) ReadURLByAlias(ctx context.Context, alias vo.Alias) (*model.URL, error) {
	const op = "postgres.url.Repository.ReadURLByAlias"

	rawUrl, err := r.queries.GetByAlias(ctx, alias.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errorx.ErrAliasDoesNotExists)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	url := converters.ConvertGenToDomain(rawUrl)
	return &url, nil
}

func (r *Repository) CreateStat(ctx context.Context, p params.CreateURLStat) (*model.URLStat, error) {
	const op = "postgres.url.Repository.CreateStat"

	rawStat, err := r.queries.CreateURLStat(ctx, converters.ConvertCreateStatParams(p))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stat := converters.ConvertGenToDomainStats(rawStat)
	return &stat, nil
}

func (r *Repository) GetCompleteAnalytics(ctx context.Context, alias vo.Alias) (*model.Analytics, error) {
	const op = "postgres.url.Repository.GetCompleteAnalytics"

	row, err := r.queries.GetCompleteAnalytics(ctx, alias.String())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	analytics, err := converters.ConvertAnalyticsToDomain(row)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return analytics, nil
}
