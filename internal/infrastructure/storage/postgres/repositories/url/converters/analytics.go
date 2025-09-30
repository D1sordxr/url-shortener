package converters

import (
	"encoding/json"
	"fmt"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/model"
	"github.com/D1sordxr/url-shortener/internal/infrastructure/storage/postgres/repositories/url/gen"
	"sort"
	"time"
)

func ConvertAnalyticsToDomain(row gen.GetCompleteAnalyticsRow) (*model.Analytics, error) {
	const op = "converters.ConvertAnalyticsToDomain"

	analytics := &model.Analytics{
		Alias:          row.Alias,
		OriginalURL:    row.OriginalUrl,
		TotalVisits:    row.TotalVisits,
		UniqueVisitors: row.UniqueVisitors,
	}

	if row.FirstVisit != nil {
		if firstVisit, ok := row.FirstVisit.(time.Time); ok {
			analytics.FirstVisit = &firstVisit
		}
	}

	if row.LastVisit != nil {
		if lastVisit, ok := row.LastVisit.(time.Time); ok {
			analytics.LastVisit = &lastVisit
		}
	}

	if err := parseRawStats(row.RawStats, analytics); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	analytics.DailyStats = aggregateDailyStats(analytics.Visits)

	analytics.UserAgentStats = aggregateUserAgentStats(analytics.Visits)

	return analytics, nil
}

func parseRawStats(rawStats json.RawMessage, analytics *model.Analytics) error {
	if len(rawStats) == 0 {
		analytics.Visits = []model.VisitStat{}
		return nil
	}

	var rawVisits []struct {
		Date      string    `json:"date"`
		UserAgent string    `json:"user_agent"`
		IPAddress string    `json:"ip_address"`
		Referer   string    `json:"referer"`
		CreatedAt time.Time `json:"created_at"`
	}

	if err := json.Unmarshal(rawStats, &rawVisits); err != nil {
		return fmt.Errorf("parse raw stats: %w", err)
	}

	visits := make([]model.VisitStat, 0, len(rawVisits))
	for _, raw := range rawVisits {
		var date time.Time
		if raw.Date != "" {
			if parsedDate, err := time.Parse("2006-01-02", raw.Date); err == nil {
				date = parsedDate
			}
		}

		visit := model.VisitStat{
			Date:      date,
			UserAgent: raw.UserAgent,
			IPAddress: raw.IPAddress,
			Referer:   raw.Referer,
			CreatedAt: raw.CreatedAt,
		}
		visits = append(visits, visit)
	}

	analytics.Visits = visits
	return nil
}

func aggregateDailyStats(visits []model.VisitStat) []model.DailyStat {
	dailyMap := make(map[time.Time]int64)

	for _, visit := range visits {
		date := time.Date(
			visit.CreatedAt.Year(),
			visit.CreatedAt.Month(),
			visit.CreatedAt.Day(),
			0,
			0,
			0,
			0,
			visit.CreatedAt.Location(),
		)
		dailyMap[date]++
	}

	dailyStats := make([]model.DailyStat, 0, len(dailyMap))
	for date, count := range dailyMap {
		dailyStats = append(dailyStats, model.DailyStat{
			Date:       date,
			VisitCount: count,
		})
	}

	sortDailyStats(dailyStats)
	return dailyStats
}

func aggregateUserAgentStats(visits []model.VisitStat) []model.UserAgentStat {
	uaMap := make(map[string]int64)

	for _, visit := range visits {
		if visit.UserAgent != "" {
			uaMap[visit.UserAgent]++
		}
	}

	uaStats := make([]model.UserAgentStat, 0, len(uaMap))
	for ua, count := range uaMap {
		uaStats = append(uaStats, model.UserAgentStat{
			UserAgent:  ua,
			VisitCount: count,
		})
	}

	return uaStats
}

func sortDailyStats(stats []model.DailyStat) {
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Date.Before(stats[j].Date)
	})
}

func sortUserAgentStats(stats []model.UserAgentStat) {
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].VisitCount < stats[j].VisitCount
	})
}
