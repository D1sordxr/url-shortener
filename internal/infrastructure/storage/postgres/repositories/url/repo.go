package url

import (
	"context"
	"fmt"
	appPorts "github.com/D1sordxr/url-shortener/internal/domain/app/ports"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/model"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/params"
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
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	url := converters.ConvertGenToDomain(rawUrl)
	return &url, nil
}

func (r *Repository) ReadByAlias(ctx context.Context, alias string) (*model.URL, error) {
	const op = "postgres.url.Repository.ReadByAlias"

	rawUrl, err := r.queries.GetByAlias(ctx, alias)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	url := converters.ConvertGenToDomain(rawUrl)
	return &url, nil
}
