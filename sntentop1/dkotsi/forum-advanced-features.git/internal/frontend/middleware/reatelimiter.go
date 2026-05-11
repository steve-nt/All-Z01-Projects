package middleware

import (
	"forum-advanced-features/internal/backend/models"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

func (i *IPRateLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println("could not read the ip address of client from the request")
		}

		limiter := i.GetLimiter(ip)
		if limiter.Allow() {

			next.ServeHTTP(w, r)
		} else {
			log.Println("client has surpassed rate limit")
		}
	})
}

type IPRateLimiter struct {
	config   *models.Config
	limiters map[string]*RateLimiter
	mutex    sync.Mutex
}

func NewIPRateLimiter(conf *models.Config) *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*RateLimiter),
		config:   conf,
	}
}
func InitializeIPRateLimiter(conf *models.Config) *IPRateLimiter {

	limiter := NewIPRateLimiter(conf)
	RefreshIPRateLimiter(limiter)

	return limiter
}

func RefreshIPRateLimiter(limiter *IPRateLimiter) {
	duration := time.Duration(time.Second) * time.Duration(limiter.config.Durations.RateLimitRefreshRate)

	ticker := time.NewTicker(duration)

	go func() {
		for range ticker.C {
			limiter.limiters = make(map[string]*RateLimiter)
		}
	}()
}
func (i *IPRateLimiter) GetLimiter(ip string) *RateLimiter {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	limiter, exists := i.limiters[ip]
	if !exists {
		//here we assign how many requests should be allowed per minute based on the const
		limiter = NewRateLimiter(i.config.Durations.RateLimit)
		i.limiters[ip] = limiter
	}
	return limiter
}

type RateLimiter struct {
	tokens         float64
	maxTokens      float64
	refillRate     float64
	lastRefillTime time.Time
	mutex          sync.Mutex
}

func NewRateLimiter(maxTokens float64) *RateLimiter {
	//maxTokens refers to how many requests per minute are allowed and the corresponding requests to achive that are added per second
	rate := maxTokens / 60.0
	return &RateLimiter{
		tokens:         maxTokens,
		maxTokens:      maxTokens,
		refillRate:     rate,
		lastRefillTime: time.Now(),
	}

}
func (l *RateLimiter) RefillTokens() {

	duration := time.Since(l.lastRefillTime).Seconds()
	tokensToAdd := duration * l.refillRate

	l.tokens += tokensToAdd
	if l.tokens > l.maxTokens {
		l.tokens = l.maxTokens
	}

	l.lastRefillTime = time.Now()

}

func (l *RateLimiter) Allow() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.RefillTokens()

	if l.tokens >= 1 {
		l.tokens--
		return true
	}

	return false
}
