package params

import (
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/vo"
	"net"
)

type CreateURL struct {
	Alias vo.Alias `json:"alias"`
	URL   vo.URL   `json:"url"`
}

type CreateURLStat struct {
	UrlID     int64       `json:"url_id"`
	UserID    *string     `json:"user_id"`
	UserAgent *string     `json:"user_agent"`
	IpAddress *net.IPAddr `json:"ip_address"`
	Referer   *string     `json:"referer"`
}
