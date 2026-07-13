package video

import (
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	interval time.Duration
	limiters *expirable.LRU[int64, *rate.Limiter]
}

func NewRateLimiter(interval time.Duration) *RateLimiter {
	return new(RateLimiter{
		interval: interval,
		limiters: expirable.NewLRU[int64, *rate.Limiter](0, nil, interval*10),
	})
}

func (r *RateLimiter) Allow(userID int64) bool {
	limiter, ok := r.limiters.Get(userID)
	if !ok {
		limiter = rate.NewLimiter(rate.Every(r.interval), 1)
		r.limiters.Add(userID, limiter)
	}

	return limiter.Allow()
}
