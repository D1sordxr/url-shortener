package converters

import (
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/model"
	"github.com/D1sordxr/url-shortener/internal/transport/http/api/url/handler/gen"
	"github.com/oapi-codegen/runtime/types"
)

func ConvertModelToResponse(m *model.URL) gen.URLResponse {
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

func ConvertAnalyticsToResponse(analytics *model.Analytics) gen.AnalyticsResponse {
	response := gen.AnalyticsResponse{
		Alias:          analytics.Alias,
		OriginalUrl:    analytics.OriginalURL,
		TotalVisits:    analytics.TotalVisits,
		UniqueVisitors: analytics.UniqueVisitors,
		Visits:         convertVisitsToResponse(analytics.Visits),
		DailyStats:     convertDailyStatsToResponse(analytics.DailyStats),
		UserAgentStats: convertUserAgentStatsToResponse(analytics.UserAgentStats),
	}

	if analytics.FirstVisit != nil {
		firstVisit := *analytics.FirstVisit
		response.FirstVisit = &firstVisit
	}
	if analytics.LastVisit != nil {
		lastVisit := *analytics.LastVisit
		response.LastVisit = &lastVisit
	}

	return response
}

func convertVisitsToResponse(visits []model.VisitStat) []gen.VisitStat {
	result := make([]gen.VisitStat, 0, len(visits))
	for _, visit := range visits {
		result = append(result, gen.VisitStat{
			Date:      types.Date{Time: visit.Date},
			UserAgent: visit.UserAgent,
			IpAddress: visit.IPAddress,
			Referer:   &visit.Referer,
			CreatedAt: visit.CreatedAt,
		})
	}
	return result
}

func convertDailyStatsToResponse(stats []model.DailyStat) []gen.DailyStat {
	result := make([]gen.DailyStat, 0, len(stats))
	for _, stat := range stats {
		result = append(result, gen.DailyStat{
			Date:       types.Date{Time: stat.Date},
			VisitCount: stat.VisitCount,
		})
	}
	return result
}

func convertUserAgentStatsToResponse(stats []model.UserAgentStat) []gen.UserAgentStat {
	result := make([]gen.UserAgentStat, 0, len(stats))
	for _, stat := range stats {
		result = append(result, gen.UserAgentStat{
			UserAgent:  stat.UserAgent,
			VisitCount: stat.VisitCount,
		})
	}
	return result
}
