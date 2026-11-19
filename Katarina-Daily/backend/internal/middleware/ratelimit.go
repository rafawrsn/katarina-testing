package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"fashion-ecommerce-backend/internal/utils"
)

// visitor tracks a single client IP's token bucket limiter.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// ipRateLimiter hands out (and periodically garbage-collects) a
// rate.Limiter per client IP.
type ipRateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rps      rate.Limit
	burst    int
}

func newIPRateLimiter(requestsPerMinute int) *ipRateLimiter {
	rl := &ipRateLimiter{
		visitors: make(map[string]*visitor),
		rps:      rate.Limit(float64(requestsPerMinute) / 60.0),
		burst:    max(requestsPerMinute/4, 5),
	}
	go rl.cleanupLoop()
	return rl
}

func (rl *ipRateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.rps, rl.burst)
		rl.visitors[ip] = &visitor{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}
	v.lastSeen = time.Now()
	return v.limiter
}

func (rl *ipRateLimiter) cleanupLoop() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 5*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}


// RateLimit throttles requests per client IP to requestsPerMinute, with
// a small burst allowance. Intended for auth and payment-sensitive
// endpoints; the whole API can also be wrapped if desired.
func RateLimit(requestsPerMinute int) gin.HandlerFunc {
	limiter := newIPRateLimiter(requestsPerMinute)

	return func(c *gin.Context) {
		if !limiter.getLimiter(c.ClientIP()).Allow() {
			utils.Fail(c, http.StatusTooManyRequests, "too many requests, please slow down", "RATE_LIMITED")
			c.Abort()
			return
		}
		c.Next()
	}
}
