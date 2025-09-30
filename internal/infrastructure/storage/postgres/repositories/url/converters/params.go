package converters

import (
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/params"
	"github.com/D1sordxr/url-shortener/internal/infrastructure/storage/postgres/repositories/url/gen"
	"github.com/D1sordxr/url-shortener/pkg/sqlutil"
)

func ConvertCreateParams(p params.CreateURL) gen.CreateURLParams {
	return gen.CreateURLParams{
		Alias: p.Alias.String(),
		Url:   p.URL.String(),
	}
}

func ConvertCreateStatParams(p params.CreateURLStat) gen.CreateURLStatParams {
	return gen.CreateURLStatParams{
		UrlID:     int32(p.UrlID),
		UserID:    sqlutil.ToNullString(p.UserID),
		UserAgent: sqlutil.ToNullString(p.UserAgent),
		IpAddress: sqlutil.ToInet(p.IpAddress),
		Referer:   sqlutil.ToNullString(p.Referer),
	}
}
