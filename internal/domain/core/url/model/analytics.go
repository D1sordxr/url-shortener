package model

import "time"

type Analytics struct {
	Alias          string          `json:"alias"`
	OriginalURL    string          `json:"original_url"`
	TotalVisits    int64           `json:"total_visits"`
	UniqueVisitors int64           `json:"unique_visitors"`
	FirstVisit     *time.Time      `json:"first_visit,omitempty"`
	LastVisit      *time.Time      `json:"last_visit,omitempty"`
	Visits         []VisitStat     `json:"visits"`
	DailyStats     []DailyStat     `json:"daily_stats"`
	UserAgentStats []UserAgentStat `json:"user_agent_stats"`
}

type VisitStat struct {
	Date      time.Time `json:"date"`
	UserAgent string    `json:"user_agent"`
	IPAddress string    `json:"ip_address"`
	Referer   string    `json:"referer,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type DailyStat struct {
	Date       time.Time `json:"date"`
	VisitCount int64     `json:"visit_count"`
}

type UserAgentStat struct {
	UserAgent  string `json:"user_agent"`
	VisitCount int64  `json:"visit_count"`
}
