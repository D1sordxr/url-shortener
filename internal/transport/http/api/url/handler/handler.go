package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/D1sordxr/url-shortener/internal/application/url/input"
	"github.com/D1sordxr/url-shortener/internal/application/url/port"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/errorx"
	"github.com/D1sordxr/url-shortener/internal/transport/http/api/url/handler/converters"
	"github.com/D1sordxr/url-shortener/internal/transport/http/api/url/handler/gen"
	"github.com/D1sordxr/url-shortener/pkg/errorz"
)

type Handler struct {
	uc port.UseCase
}

func NewHandler(uc port.UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) PostShorten(
	ctx context.Context,
	request gen.PostShortenRequestObject,
) (gen.PostShortenResponseObject, error) {
	if request.Body.Url == "" {
		return gen.PostShorten400JSONResponse{Error: "url is required"}, nil
	}

	urlModel, err := h.uc.Create(ctx, input.Create{
		URL:   request.Body.Url,
		Alias: request.Body.Alias,
	})
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrAliasAlreadyExists):
			return gen.PostShorten409JSONResponse{Error: "alias already exists"}, nil
		case errorz.In(
			err,
			errorx.ErrInvalidAliasLength,
			errorx.ErrURLEmpty,
			errorx.ErrURLInvalidFormat,
			errorx.ErrURLMissingScheme,
			errorx.ErrURLMissingHost,
			errorx.ErrURLUnsupportedScheme,
			errorx.ErrURLParseFailed,
		):
			return gen.PostShorten400JSONResponse{Error: fmt.Sprintf("invalid url data: %v", err)}, nil
		default:
			return gen.PostShorten500JSONResponse{Error: fmt.Sprintf(
				"%s: %s", "internal server error", err.Error()),
			}, nil
		}
	}

	return gen.PostShorten201JSONResponse(converters.ConvertModelToResponse(urlModel)), nil
}

func (h *Handler) GetAnalyticsAlias(
	ctx context.Context,
	request gen.GetAnalyticsAliasRequestObject,
) (gen.GetAnalyticsAliasResponseObject, error) {
	if request.Alias == "" {
		return gen.GetAnalyticsAlias404JSONResponse{Error: "alias is required"}, nil
	}

	analytics, err := h.uc.GetCompleteAnalytics(ctx, request.Alias)
	if err != nil {
		switch {
		case errorz.In(err, errorx.ErrAliasDoesNotExists, errorx.ErrInvalidAliasLength):
			return gen.GetAnalyticsAlias404JSONResponse{Error: "invalid alias"}, nil
		default:
			return gen.GetAnalyticsAlias500JSONResponse{Error: "internal server error"}, nil
		}

	}

	return gen.GetAnalyticsAlias200JSONResponse(converters.ConvertAnalyticsToResponse(analytics)), nil
}

func (h *Handler) GetSAlias(
	ctx context.Context,
	request gen.GetSAliasRequestObject,
) (gen.GetSAliasResponseObject, error) {
	if request.Alias == "" {
		return gen.GetSAlias404JSONResponse{Error: "alias is required"}, nil
	}

	url, err := h.uc.GetForRedirect(ctx, input.GetForRedirect{
		Alias:     "",
		UserID:    nil,
		UserAgent: nil,
		IpAddress: nil,
		Referer:   nil,
	})
	if err != nil {
		switch {
		case errorz.In(err, errorx.ErrAliasDoesNotExists, errorx.ErrInvalidAliasLength):
			return gen.GetSAlias404JSONResponse{Error: "invalid alias"}, nil
		default:
			return gen.GetSAlias500JSONResponse{Error: "internal server error"}, nil
		}
	}

	return gen.GetSAlias302Response{
		Headers: gen.GetSAlias302ResponseHeaders{Location: url.String()},
	}, nil
}
