package converters

import (
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/model"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/vo"
	"github.com/D1sordxr/url-shortener/internal/infrastructure/storage/postgres/repositories/url/gen"
)

func ConvertGenToDomain(rawModel gen.Url) model.URL {
	return model.URL{
		ID:        int64(rawModel.ID),
		Alias:     vo.Alias(rawModel.Alias),
		URL:       vo.URL(rawModel.Url),
		CreatedAt: rawModel.CreatedAt,
		UpdatedAt: rawModel.UpdatedAt,
	}
}

func ConvertGenToDomainStats(rawStat gen.UrlStat) model.URLStat {
	return model.URLStat{
		ID:        int64(rawStat.ID),
		UrlID:     int64(rawStat.UrlID),
		UserID:    rawStat.UserID.String,
		UserAgent: rawStat.UserAgent.String,
		IPAddress: rawStat.IpAddress.IPNet.String(),
		Referer:   rawStat.Referer.String,
		CreatedAt: rawStat.CreatedAt,
	}
}

func ConvertGenToDomainRowStats(rawStat gen.GetUrlStatsRow) model.URLStat {
	return model.URLStat{
		ID:        int64(rawStat.ID),
		UrlID:     int64(rawStat.UrlID),
		Alias:     rawStat.Alias,
		Url:       rawStat.Url,
		UserID:    rawStat.UserID.String,
		UserAgent: rawStat.UserAgent.String,
		IPAddress: rawStat.IpAddress.IPNet.String(),
		Referer:   rawStat.Referer.String,
		CreatedAt: rawStat.CreatedAt,
	}
}
