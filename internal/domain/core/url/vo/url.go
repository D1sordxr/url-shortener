package vo

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/D1sordxr/url-shortener/internal/domain/core/url/errorx"
)

type URL string

const (
	safeUrlPrefix = "https://"
)

func NewURL(rawURL string) (URL, error) {
	if err := validateURL(rawURL); err != nil {
		return "", err
	}

	normalizedURL := normalizeURL(rawURL)

	return URL(normalizedURL), nil
}

func Validate[T comparable](url T) bool {
	var urlStr string

	switch u := any(url).(type) {
	case string:
		urlStr = u
	case URL:
		urlStr = string(u)
	default:
		return false
	}

	return validateURL(urlStr) == nil
}

func (u URL) String() string {
	return string(u)
}

func (u URL) IsSecure() bool {
	return strings.HasPrefix(string(u), safeUrlPrefix)
}

func (u URL) GetDomain() string {
	parsed, err := url.Parse(string(u))
	if err != nil {
		return ""
	}
	return parsed.Host
}

func (u URL) WithScheme(scheme string) (URL, error) {
	parsed, err := url.Parse(string(u))
	if err != nil {
		return "", fmt.Errorf("%w: %w", errorx.ErrURLParseFailed, err)
	}

	parsed.Scheme = scheme
	return URL(parsed.String()), nil
}

func (u URL) Equals(other URL) bool {
	return u.Normalize() == other.Normalize()
}

func (u URL) Normalize() URL {
	parsed, err := url.Parse(string(u))
	if err != nil {
		return u
	}

	parsed.Path = strings.TrimSuffix(parsed.Path, "/")
	if parsed.Path == "" {
		parsed.Path = "/"
	}

	return URL(parsed.String())
}

func validateURL(rawURL string) error {
	if strings.TrimSpace(rawURL) == "" {
		return errorx.ErrURLEmpty
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("%w: %w", errorx.ErrURLInvalidFormat, err)
	}

	if parsed.Scheme == "" {
		return errorx.ErrURLMissingScheme
	}

	if parsed.Host == "" {
		return errorx.ErrURLMissingHost
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errorx.ErrURLUnsupportedScheme
	}

	return nil
}

func normalizeURL(rawURL string) string {
	if !strings.Contains(rawURL, "://") {
		return safeUrlPrefix + rawURL
	}
	return rawURL
}
