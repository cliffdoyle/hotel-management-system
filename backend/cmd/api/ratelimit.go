package main

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// rateLimit checks if the request from a given IP exceeds the limit.
func (app *application) withRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// These values should be configurable, perhaps from .env.
		// For example, 100 requests per 1 minute window.
		const requestsPerWindow = 100
		const window = 1 * time.Minute

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		app.logger.Info("request-received","ip",ip)

		key := fmt.Sprintf("rate_limit:%s", ip)
		now := time.Now()
		windowStart := now.Add(-window).UnixNano()

		ctx := r.Context()
		
		// Use a transaction to ensure atomicity.
		tx := app.redis.TxPipeline()
		tx.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))
		tx.ZAdd(ctx, key, &redis.Z{Score: float64(now.UnixNano()), Member: now.UnixNano()})
		count := tx.ZCard(ctx, key)
		tx.Expire(ctx, key, window) // Set an expiry to clean up old keys.
		
		_, err = tx.Exec(ctx)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		
		currentCount := count.Val()

		// Set response headers
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(requestsPerWindow))
		w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(int64(requestsPerWindow)-currentCount, 10))
		resetTime := now.Add(window).Unix()
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

		// If limit is exceeded, send a 429 response.
		if currentCount > requestsPerWindow {
			message := "rate limit exceeded"
			app.errorResponse(w, r, http.StatusTooManyRequests, message)
			return
		}

		next.ServeHTTP(w, r)
	})
}