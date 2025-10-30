package headerctx

import (
	"context"

	"github.com/labstack/echo/v4"
)

type HeaderCtxKey string

// InjectHeaders is a middleware that injects the given headers into echo.Context and context.Context.
// A boolean 'require' dictates whether all headers need to be present for request to continue.
// If any headers are missing the middleware returns a 500 Internal server error.
func InjectHeaders(require bool, headers ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			ctx := req.Context()
			for _, key := range headers {
				if value := req.Header.Get(key); value != "" {
					c.Set(key, value)
					ctx = context.WithValue(ctx, HeaderCtxKey(key), value)
					continue
				}
				// return 500 Internal server error if headers are required, but not present
				if require {
					return echo.NewHTTPError(500, "Internal server error")
				}
			}
			c.SetRequest(req.WithContext(ctx))
			return next(c)
		}
	}
}
