package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/D1sordxr/url-shortener/internal/application/url/input"
	"github.com/D1sordxr/url-shortener/internal/application/url/port"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/errorx"
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
	request PostShortenRequestObject,
) (PostShortenResponseObject, error) {
	if request.Body.Url == "" {
		return PostShorten400JSONResponse{
			Error: "url is required",
		}, nil
	}

	urlModel, err := h.uc.Create(ctx, input.Create{
		URL:   request.Body.Url,
		Alias: request.Body.Alias,
	})
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrAliasAlreadyExists):
			return PostShorten409JSONResponse{
				Error: errorx.ErrAliasAlreadyExists.Error(),
			}, nil
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
			return PostShorten400JSONResponse{
				Error: fmt.Sprintf("invalid url data: %v", err),
			}, nil
		default:
			return PostShorten500JSONResponse{
				Error: fmt.Sprintf("%s: %s", "internal server error", err.Error()),
			}, nil
		}
	}

	resp := ConvertModelToUrlResponse(urlModel)
	return PostShorten201JSONResponse(resp), nil
}

func (h *Handler) GetAnalyticsAlias(
	ctx context.Context,
	request GetAnalyticsAliasRequestObject,
) (GetAnalyticsAliasResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) GetSAlias(
	ctx context.Context,
	request GetSAliasRequestObject,
) (GetSAliasResponseObject, error) {
	//TODO implement me
	panic("implement me")
}
