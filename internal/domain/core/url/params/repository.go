package params

import "github.com/D1sordxr/url-shortener/internal/domain/core/url/vo"

type CreateURL struct {
	Alias vo.Alias `json:"alias"`
	URL   vo.URL   `json:"url"`
}
