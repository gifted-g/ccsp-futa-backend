package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type bucket struct { tokens int; last time.Time }

// Simple IP-based token bucket: limit reqsPerMin/min
func RateLimit(reqsPerMin int) gin.HandlerFunc {
	var mu sync.Mutex
	m := map[string]*bucket{}
	refill := func(b *bucket){
		d := time.Since(b.last).Minutes()
		b.tokens += int(d*float64(reqsPerMin))
		if b.tokens > reqsPerMin { b.tokens = reqsPerMin }
		b.last = time.Now()
	}
	return func(c *gin.Context) {
		ip, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
		mu.Lock()
		b := m[ip]
		if b == nil { b = &bucket{tokens: reqsPerMin, last: time.Now()}; m[ip]=b }
		refill(b)
		if b.tokens <= 0 { mu.Unlock(); c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error":"rate limit"}); return }
		b.tokens--
		mu.Unlock()
		c.Next()
	}
}