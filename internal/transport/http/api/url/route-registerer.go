package url

import (
	"github.com/D1sordxr/url-shortener/internal/transport/http/api/url/handler"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

type RouteRegisterer struct {
	handler     *handler.Handler
	middlewares []gin.HandlerFunc
}

func NewRouteRegisterer(
	handlers *handler.Handler,
	middlewares ...gin.HandlerFunc,
) *RouteRegisterer {
	return &RouteRegisterer{
		handler:     handlers,
		middlewares: middlewares,
	}
}

func (r *RouteRegisterer) RegisterRoutes(router *ginext.RouterGroup) {
	urlGroup := router.RouterGroup
	for _, mw := range r.middlewares {
		urlGroup.Use(mw)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, map[string]any{"message": "OK"})
	})

	handler.RegisterHandlers(
		urlGroup,
		handler.NewStrictHandler(
			r.handler,
			[]handler.StrictMiddlewareFunc{},
		),
	)
}
