package model

import (
	"time"
)

type URL struct {
	ID        int64     `json:"id" db:"id"`
	Alias     string    `json:"alias" db:"alias"`
	URL       string    `json:"url" db:"url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type URLStat struct {
	ID        int64     `json:"id" db:"id"`
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
