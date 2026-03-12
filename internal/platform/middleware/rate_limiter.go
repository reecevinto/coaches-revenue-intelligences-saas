package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type visitor struct {
	lastSeen time.Time
	count    int
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

// configuration
const (
	requestLimit = 5           // max requests
	window       = time.Second // per time window
	cleanupTime  = time.Minute * 2
)

// RateLimiter middleware protects against request flooding.
func RateLimiter(next http.Handler) http.Handler {

	// start background cleanup goroutine
	go cleanupVisitors()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		mu.Lock()

		v, exists := visitors[ip]

		if !exists {
			visitors[ip] = &visitor{
				lastSeen: time.Now(),
				count:    1,
			}
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		// reset window if expired
		if time.Since(v.lastSeen) > window {
			v.count = 1
			v.lastSeen = time.Now()
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		// increment count
		v.count++

		if v.count > requestLimit {
			mu.Unlock()

			http.Error(
				w,
				"too many requests",
				http.StatusTooManyRequests,
			)
			return
		}

		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

// cleanupVisitors removes stale IP entries
func cleanupVisitors() {

	for {

		time.Sleep(cleanupTime)

		mu.Lock()

		for ip, v := range visitors {

			if time.Since(v.lastSeen) > cleanupTime {
				delete(visitors, ip)
			}

		}

		mu.Unlock()
	}
}
