package converters

import (
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/model"
	"github.com/D1sordxr/url-shortener/internal/infrastructure/storage/postgres/repositories/url/gen"
)

func ConvertGenToDomain(rawModel gen.Url) model.URL {
	return model.URL{
		ID:        int64(rawModel.ID),
		Alias:     rawModel.Alias,
		URL:       rawModel.Url,
		CreatedAt: rawModel.CreatedAt,
		UpdatedAt: rawModel.UpdatedAt,
	}
}
