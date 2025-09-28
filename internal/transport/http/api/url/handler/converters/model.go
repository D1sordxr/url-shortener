package converters

import (
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/model"
	"github.com/D1sordxr/url-shortener/internal/transport/http/api/url/handler/gen"
)

func ConvertModelToUrlResponse(m *model.URL) gen.URLResponse {
	alias := m.Alias.String()
	originalUrl := m.URL.String()
	shortUrl := "/api/url/s/" + alias
	return gen.URLResponse{
		Alias:       &alias,
		CreatedAt:   &m.CreatedAt,
		Id:          &m.ID,
		OriginalUrl: &originalUrl,
		ShortUrl:    &shortUrl,
	}
}
