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
	"net"
	"time"
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

func (uc *UseCase) GetForRedirect(ctx context.Context, i input.GetForRedirect) (vo.URL, error) {
	const op = "url.UseCase.GetForRedirect"
	logFields := logger.WithFields("alias", i.Alias)

	uc.log.Info("Attempting to get URL by alias", logFields()...)

	aliasValue, err := vo.NewAlias(i.Alias)
	if err != nil {
		uc.log.Error("Error parsing URL alias", logFields("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	url, err := uc.repo.ReadURLByAlias(ctx, aliasValue)
	if err != nil {
		uc.log.Error("Error getting URL by alias", logFields("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	go func() {
		statCtx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		var ip net.IP
		if i.IpAddress != nil {
			ip = net.ParseIP(*i.IpAddress)
		}

		stat, statErr := uc.repo.CreateStat(statCtx, params.CreateURLStat{
			UrlID:     url.ID,
			UserID:    i.UserID,
			UserAgent: i.UserAgent,
			IpAddress: &ip,
			Referer:   i.Referer,
		})
		if statErr != nil {
			uc.log.Error("Error creating stat", logFields("error", statErr.Error()))
			return
		}
		uc.log.Info("URL stat created", logFields(
			"url_id", url.ID,
			"alias", i.Alias,
			"user_id", i.UserID,
			"user_agent", i.UserAgent,
			"ip_address", i.IpAddress,
			"referer", i.Referer,
			"stat", stat,
		)...)
	}()

	uc.log.Info("URL for redirect retrieved successfully", logFields()...)

	return url.URL, nil
}

func (uc *UseCase) GetCompleteAnalytics(ctx context.Context, alias string) (*model.Analytics, error) {
	const op = "url.UseCase.GetCompleteAnalytics"
	logFields := logger.WithFields("operation", op, "alias", alias)

	uc.log.Info("Attempting to get analytics by alias", logFields()...)

	validAlias, err := vo.NewAlias(alias)
	if err != nil {
		uc.log.Error("Error parsing URL alias", logFields("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	analytics, err := uc.repo.GetCompleteAnalytics(ctx, validAlias)
	if err != nil {
		uc.log.Error("Error getting analytics by alias", logFields("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	uc.log.Info("Successfully got analytics by alias", logFields()...)

	return analytics, nil
}
