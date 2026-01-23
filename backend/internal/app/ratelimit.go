package app

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimitConfig defines rate limiting parameters for an endpoint or group.
type RateLimitConfig struct {
	// Rate is the number of requests allowed per second.
	Rate rate.Limit
	// Burst is the maximum number of requests allowed in a burst.
	Burst int
	// KeyFunc extracts the rate limit key from a request.
	// If nil, defaults to IP-based limiting.
	KeyFunc func(r *http.Request) string
}

// Common rate limit presets
var (
	// AuthStrict: 5 requests per minute (for login/register)
	AuthStrict = RateLimitConfig{
		Rate:  rate.Limit(5.0 / 60.0), // 5 per minute
		Burst: 5,
	}

	// RegisterStrict: 5 requests per hour (account creation)
	RegisterStrict = RateLimitConfig{
		Rate:  rate.Limit(5.0 / 3600.0), // 5 per hour
		Burst: 3,
	}

	// APIGeneral: 120 requests per minute (authenticated endpoints)
	APIGeneral = RateLimitConfig{
		Rate:  rate.Limit(120.0 / 60.0), // 120 per minute = 2/sec
		Burst: 20,
	}
)

// RateLimiter manages rate limiters keyed by some identifier (IP, user ID, etc).
// NOTE: In-memory best-effort only; limits reset on restart and are per-instance.
type RateLimiter struct {
	limiters        map[string]*limiterEntry
	mu              sync.RWMutex
	config          RateLimitConfig
	cleanupInterval time.Duration
	entryTTL        time.Duration
	stopCleanup     chan struct{}
}

type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewRateLimiter creates a new rate limiter with the given configuration.
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	entryTTL := 10 * time.Minute
	if config.Rate > 0 && config.Burst > 0 {
		refill := time.Duration(float64(config.Burst)/float64(config.Rate)) * time.Second
		if refill > 0 && refill*2 > entryTTL {
			entryTTL = refill * 2
		}
	}

	rl := &RateLimiter{
		limiters:        make(map[string]*limiterEntry),
		config:          config,
		cleanupInterval: 5 * time.Minute,
		entryTTL:        entryTTL,
		stopCleanup:     make(chan struct{}),
	}
	go rl.cleanupLoop()
	return rl
}

// Allow checks if a request with the given key is allowed.
// Returns true if allowed, false if rate limited.
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	entry, exists := rl.limiters[key]
	if !exists {
		entry = &limiterEntry{
			limiter:  rate.NewLimiter(rl.config.Rate, rl.config.Burst),
			lastSeen: time.Now(),
		}
		rl.limiters[key] = entry
	} else {
		entry.lastSeen = time.Now()
	}
	rl.mu.Unlock()

	return entry.limiter.Allow()
}

// Reserve returns a reservation for the given key, allowing callers to
// determine how long to wait before the request would be allowed.
func (rl *RateLimiter) Reserve(key string) *rate.Reservation {
	rl.mu.Lock()
	entry, exists := rl.limiters[key]
	if !exists {
		entry = &limiterEntry{
			limiter:  rate.NewLimiter(rl.config.Rate, rl.config.Burst),
			lastSeen: time.Now(),
		}
		rl.limiters[key] = entry
	} else {
		entry.lastSeen = time.Now()
	}
	rl.mu.Unlock()

	return entry.limiter.Reserve()
}

// cleanupLoop periodically removes stale limiters to prevent memory leaks.
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.cleanup()
		case <-rl.stopCleanup:
			return
		}
	}
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	cutoff := time.Now().Add(-rl.entryTTL)
	for key, entry := range rl.limiters {
		if entry.lastSeen.Before(cutoff) {
			delete(rl.limiters, key)
		}
	}
}

// Stop halts the cleanup goroutine.
func (rl *RateLimiter) Stop() {
	close(rl.stopCleanup)
}

// RateLimiters holds all rate limiters used by the application.
type RateLimiters struct {
	// Auth limits login/register by IP
	Auth *RateLimiter
	// Register has stricter limits for account creation
	Register *RateLimiter
	// API limits authenticated endpoints by user ID
	API *RateLimiter
}

// NewRateLimiters creates all rate limiters with default configurations.
func NewRateLimiters() *RateLimiters {
	return &RateLimiters{
		Auth:     NewRateLimiter(AuthStrict),
		Register: NewRateLimiter(RegisterStrict),
		API:      NewRateLimiter(APIGeneral),
	}
}

// NewTestRateLimiters creates rate limiters with permissive limits for testing.
func NewTestRateLimiters() *RateLimiters {
	permissive := RateLimitConfig{
		Rate:  rate.Limit(1000), // 1000/sec
		Burst: 1000,
	}
	return &RateLimiters{
		Auth:     NewRateLimiter(permissive),
		Register: NewRateLimiter(permissive),
		API:      NewRateLimiter(permissive),
	}
}

// Stop halts all cleanup goroutines.
func (rl *RateLimiters) Stop() {
	rl.Auth.Stop()
	rl.Register.Stop()
	rl.API.Stop()
}

// RateLimitByIP creates middleware that rate limits by client IP.
func (app *Application) RateLimitByIP(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := app.GetClientIP(r)

			reservation := limiter.Reserve(key)
			if !reservation.OK() {
				app.respondRateLimited(w)
				return
			}

			delay := reservation.Delay()
			if delay > 0 {
				reservation.Cancel()
				app.respondRateLimitedWithRetry(w, delay, limiter.config.Burst)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitByUser creates middleware that rate limits by user ID.
// Falls back to IP if user is not authenticated.
func (app *Application) RateLimitByUser(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var key string

			// Try to get user from context (set by Authenticate middleware)
			user := app.GetUser(r.Context())
			if user != nil && !user.IsAnonymous() {
				key = "user:" + user.ID.String()
			} else {
				key = "ip:" + app.GetClientIP(r)
			}

			reservation := limiter.Reserve(key)
			if !reservation.OK() {
				app.respondRateLimited(w)
				return
			}

			delay := reservation.Delay()
			if delay > 0 {
				reservation.Cancel()
				app.respondRateLimitedWithRetry(w, delay, limiter.config.Burst)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *Application) respondRateLimited(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)
	w.Write([]byte(`{"error": "rate limit exceeded"}`))
}

func (app *Application) respondRateLimitedWithRetry(w http.ResponseWriter, delay time.Duration, limit int) {
	retryAfter := int(delay.Seconds())
	if retryAfter < 1 {
		retryAfter = 1
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
	w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
	w.Header().Set("X-RateLimit-Remaining", "0")
	w.WriteHeader(http.StatusTooManyRequests)
	w.Write([]byte(`{"error": "rate limit exceeded"}`))
}
