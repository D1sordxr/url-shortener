package model

import (
	"time"

	"github.com/D1sordxr/url-shortener/internal/domain/core/url/vo"
)

type URL struct {
	ID        int64     `json:"id" db:"id"`
	URL       vo.URL    `json:"url" db:"url"`
	Alias     vo.Alias  `json:"alias" db:"alias"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type URLStat struct {
	ID        int64     `json:"id" db:"id"`
	Alias     string    `json:"alias" db:"alias"`
	Url       string    `json:"url" db:"url"`
	UrlID     int64     `json:"url_id" db:"url_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	Referer   string    `json:"referer" db:"referer"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type URLWithStats struct {
	URL            *URL  `json:"url"`
	TotalVisits    int64 `json:"total_visits"`
	UniqueVisitors int64 `json:"unique_visitors"`
}
