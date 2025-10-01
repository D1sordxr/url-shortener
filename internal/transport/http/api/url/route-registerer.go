package url

import (
	"github.com/D1sordxr/url-shortener/internal/transport/http/api/url/handler"
	"github.com/D1sordxr/url-shortener/internal/transport/http/api/url/handler/gen"
	"github.com/wb-go/wbf/ginext"
)

type RouteRegisterer struct {
	handler     *handler.Handler
	middlewares []ginext.HandlerFunc
}

func NewRouteRegisterer(
	handlers *handler.Handler,
	middlewares ...ginext.HandlerFunc,
) *RouteRegisterer {
	return &RouteRegisterer{
		handler:     handlers,
		middlewares: middlewares,
	}
}

func (r *RouteRegisterer) RegisterRoutes(router *ginext.RouterGroup) {
	urlGroup := router.Group("/url")
	for _, mw := range r.middlewares {
		urlGroup.Use(mw)
	}

	router.GET("/health", func(c *ginext.Context) {
		c.JSON(200, map[string]any{"message": "OK"})
	})

	gen.RegisterHandlers(urlGroup, r.handler)
}
