package converters

import (
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/params"
	"github.com/D1sordxr/url-shortener/internal/infrastructure/storage/postgres/repositories/url/gen"
)

func ConvertCreateParams(p params.CreateURL) gen.CreateURLParams {
	return gen.CreateURLParams{
		Alias: p.Alias,
		Url:   p.Url,
	}
}
