package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(5)

	t.Run("allows requests under limit", func(t *testing.T) {
		ip := "192.168.1.1"
		for i := 0; i < 5; i++ {
			if !limiter.Allow(ip) {
				t.Errorf("request %d should be allowed", i+1)
			}
		}
	})

	t.Run("blocks requests over limit", func(t *testing.T) {
		ip := "192.168.1.2"
		for i := 0; i < 5; i++ {
			limiter.Allow(ip)
		}

		if limiter.Allow(ip) {
			t.Error("request should be blocked after limit")
		}
	})

	t.Run("resets after window", func(t *testing.T) {
		limiter := NewRateLimiter(2)
		ip := "192.168.1.3"

		limiter.Allow(ip)
		limiter.Allow(ip)

		if limiter.Allow(ip) {
			t.Error("should be blocked")
		}

		time.Sleep(61 * time.Second)

		if !limiter.Allow(ip) {
			t.Error("should be allowed after window reset")
		}
	})
}

func TestRateLimitMiddleware(t *testing.T) {
	limiter := NewRateLimiter(2)
	middleware := RateLimitMiddleware(limiter)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	wrappedHandler := middleware(handler)

	t.Run("allows requests under limit", func(t *testing.T) {
		for i := 0; i < 2; i++ {
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.100:1234"
			w := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("request %d: expected status 200, got %d", i+1, w.Code)
			}
		}
	})

	t.Run("blocks requests over limit", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.100:1234"
		w := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(w, req)

		if w.Code != http.StatusTooManyRequests {
			t.Errorf("expected status 429, got %d", w.Code)
		}
	})

	t.Run("uses X-Forwarded-For header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.101:1234"
		req.Header.Set("X-Forwarded-For", "10.0.0.1")
		w := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}
