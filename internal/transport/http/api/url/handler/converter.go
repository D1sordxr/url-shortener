package handler

import (
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/model"
)

func ConvertModelToUrlResponse(m *model.URL) URLResponse {
	alias := m.Alias.String()
	originalUrl := m.URL.String()
	shortUrl := "/api/url/s/" + alias
	return URLResponse{
		Alias:       &alias,
		CreatedAt:   &m.CreatedAt,
		Id:          &m.ID,
		OriginalUrl: &originalUrl,
		ShortUrl:    &shortUrl,
	}
}
