package middleware

import (
	"github.com/wb-go/wbf/ginext"
	"net/http"
	"strings"
)

const (
	ipKey        = "ip_address"
	userAgentKey = "user_agent"
	refererKey   = "referer"
	userIDKey    = "user_id"
)

var Stat ginext.HandlerFunc = func(c *ginext.Context) {
	defer c.Next()

	if c.Request.Method != http.MethodGet {
		return
	}

	ip := c.ClientIP()
	userAgent := strings.TrimSpace(c.Request.UserAgent())
	referer := strings.TrimSpace(c.Request.Referer())
	userID := "user-id-from-jwt"

	if ip != "" {
		c.Set(ipKey, ip)
	}
	if userAgent != "" {
		c.Set(userAgentKey, userAgent)
	}
	if referer != "" {
		c.Set(refererKey, referer)
	}
	c.Set(userIDKey, userID)

	c.Next()
}

func ReadIPFromCtx(c *ginext.Context) string {
	if ip, exists := c.Get(ipKey); exists {
		if ipStr, ok := ip.(string); ok {
			return ipStr
		}
	}
	return ""
}

func ReadUserAgentFromCtx(c *ginext.Context) string {
	if ua, exists := c.Get(userAgentKey); exists {
		if uaStr, ok := ua.(string); ok {
			return uaStr
		}
	}
	return ""
}

func ReadRefererFromCtx(c *ginext.Context) string {
	if ref, exists := c.Get(refererKey); exists {
		if refStr, ok := ref.(string); ok {
			return refStr
		}
	}
	return ""
}

func ReadUserIDFromCtx(c *ginext.Context) string {
	if userID, exists := c.Get(userIDKey); exists {
		if userIDStr, ok := userID.(string); ok {
			return userIDStr
		}
	}
	return ""
}
