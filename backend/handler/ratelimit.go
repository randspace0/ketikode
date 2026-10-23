package handler

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// RateLimit allows perSecond requests per client IP with a burst, and answers 429 beyond that.
// State is in memory, so each process counts separately.
func RateLimit(perSecond float64, burst int) echo.MiddlewareFunc {
	return middleware.RateLimiter(middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
		Rate:      rate.Limit(perSecond),
		Burst:     burst,
		ExpiresIn: 10 * time.Minute,
	}))
}
