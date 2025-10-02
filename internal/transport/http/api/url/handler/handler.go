package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/D1sordxr/url-shortener/internal/application/url/input"
	"github.com/D1sordxr/url-shortener/internal/application/url/port"
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/errorx"
	"github.com/D1sordxr/url-shortener/internal/transport/http/api/url/handler/converters"
	"github.com/D1sordxr/url-shortener/internal/transport/http/api/url/handler/gen"
	"github.com/D1sordxr/url-shortener/internal/transport/http/middleware"
	"github.com/D1sordxr/url-shortener/pkg/errorz"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc port.UseCase
}

func New(uc port.UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) PostShorten(c *gin.Context) {
	var request struct {
		URL   string  `json:"url" binding:"required,url"`
		Alias *string `json:"alias,omitempty" binding:"omitempty,max=128"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gen.ErrorResponse{
			Error: fmt.Sprintf("%s: %v", "invalid request body", err),
		})
		return
	}

	urlModel, err := h.uc.Create(c.Request.Context(), input.Create{
		URL:   request.URL,
		Alias: request.Alias,
	})
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrAliasAlreadyExists):
			c.JSON(http.StatusConflict, gen.ErrorResponse{
				Error: "alias already exists",
			})
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
			c.JSON(http.StatusBadRequest, gen.ErrorResponse{
				Error: fmt.Sprintf("invalid url data: %v", err),
			})
		default:
			c.JSON(http.StatusInternalServerError, gen.ErrorResponse{
				Error: "internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, converters.ConvertModelToResponse(urlModel))
}

func (h *Handler) GetAnalyticsAlias(c *gin.Context, alias string) {
	if alias == "" {
		c.JSON(http.StatusBadRequest, gen.ErrorResponse{
			Error: "alias is required",
		})
		return
	}

	analytics, err := h.uc.GetCompleteAnalytics(c.Request.Context(), alias)
	if err != nil {
		switch {
		case errorz.In(err, errorx.ErrAliasDoesNotExists, errorx.ErrInvalidAliasLength):
			c.JSON(http.StatusNotFound, gen.ErrorResponse{
				Error: "alias not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gen.ErrorResponse{
				Error: "internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, converters.ConvertAnalyticsToResponse(analytics))
}

func (h *Handler) GetSAlias(c *gin.Context, alias string) {
	if alias == "" {
		c.JSON(http.StatusBadRequest, gen.ErrorResponse{
			Error: "alias is required",
		})
		return
	}

	ip := middleware.ReadIPFromCtx(c)
	userID := middleware.ReadUserIDFromCtx(c)
	userAgent := middleware.ReadUserAgentFromCtx(c)
	referer := middleware.ReadRefererFromCtx(c)

	redirectInput := input.GetForRedirect{
		Alias:     alias,
		UserID:    &userID,
		UserAgent: &userAgent,
		IpAddress: &ip,
		Referer:   &referer,
	}

	url, err := h.uc.GetForRedirect(c.Request.Context(), redirectInput)
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrAliasDoesNotExists),
			errorz.In(err, errorx.ErrInvalidAliasLength):
			c.JSON(http.StatusNotFound, gen.ErrorResponse{
				Error: "alias not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gen.ErrorResponse{
				Error: "internal server error",
			})
		}
		return
	}

	c.Redirect(http.StatusFound, url.String())
}
