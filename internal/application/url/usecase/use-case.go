package usecase

import (
	"context"
	"fmt"
	"github.com/D1sordxr/url-shortener/internal/application/url/input"
	appPorts "github.com/D1sordxr/url-shortener/internal/domain/app/ports"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/model"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/params"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/ports"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/vo"
	"github.com/D1sordxr/url-shortener/pkg/logger"
)

type UseCase struct {
	log  appPorts.Logger
	repo ports.Repository
}

func NewUseCase(
	log appPorts.Logger,
	repo ports.Repository,
) *UseCase {
	return &UseCase{
		log:  log,
		repo: repo,
	}
}

func (uc *UseCase) Create(ctx context.Context, create input.Create) (*model.URL, error) {
	const op = "url.UseCase.Create"
	logFields := logger.WithFields("operation", op, "url", create.URL)

	uc.log.Info("Attempting to create new URL", logFields()...)

	urlValue, err := vo.NewURL(create.URL)
	if err != nil {
		uc.log.Error("Error parsing URL", logFields("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var aliasValue vo.Alias
	switch create.Alias {
	case nil:
		aliasValue = vo.GenerateAlias()
	default:
		aliasValue, err = vo.NewAlias(*create.Alias)
		if err != nil {
			uc.log.Error("Error parsing URL alias", logFields("error", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	urlModel, err := uc.repo.Create(ctx, params.CreateURL{
		Alias: aliasValue,
		URL:   urlValue,
	})
	if err != nil {
		uc.log.Error("Error creating URL", logFields("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	uc.log.Info("URL created successfully", logFields()...)

	return urlModel, nil
}
