package app

import (
	"net/http/httptest"
	"testing"
)

func TestTrustedProxyChecker(t *testing.T) {
	tests := []struct {
		name     string
		trusted  []string
		testIP   string
		expected bool
	}{
		{
			name:     "exact IP match",
			trusted:  []string{"127.0.0.1"},
			testIP:   "127.0.0.1",
			expected: true,
		},
		{
			name:     "no match",
			trusted:  []string{"127.0.0.1"},
			testIP:   "192.168.1.1",
			expected: false,
		},
		{
			name:     "CIDR match",
			trusted:  []string{"10.0.0.0/8"},
			testIP:   "10.1.2.3",
			expected: true,
		},
		{
			name:     "CIDR no match",
			trusted:  []string{"10.0.0.0/8"},
			testIP:   "192.168.1.1",
			expected: false,
		},
		{
			name:     "empty list",
			trusted:  []string{},
			testIP:   "127.0.0.1",
			expected: false,
		},
		{
			name:     "IPv6 match",
			trusted:  []string{"::1"},
			testIP:   "::1",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker, err := NewTrustedProxyChecker(tt.trusted)
			if err != nil {
				t.Fatalf("failed to create checker: %v", err)
			}

			result := checker.IsTrusted(tt.testIP)
			if result != tt.expected {
				t.Errorf("IsTrusted(%q) = %v, want %v", tt.testIP, result, tt.expected)
			}
		})
	}
}

func TestExtractIPFromAddr(t *testing.T) {
	tests := []struct {
		addr     string
		expected string
	}{
		{"127.0.0.1:8080", "127.0.0.1"},
		{"192.168.1.1:3000", "192.168.1.1"},
		{"[::1]:8080", "::1"},
		{"[2001:db8::1]:443", "2001:db8::1"},
		{"127.0.0.1", "127.0.0.1"},
		{"::1", "::1"},
	}

	for _, tt := range tests {
		t.Run(tt.addr, func(t *testing.T) {
			result := extractIPFromAddr(tt.addr)
			if result != tt.expected {
				t.Errorf("extractIPFromAddr(%q) = %q, want %q", tt.addr, result, tt.expected)
			}
		})
	}
}

func TestGetClientIP(t *testing.T) {
	app := &Application{
		trustedProxies: nil,
	}

	t.Run("no trusted proxies uses RemoteAddr", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "192.168.1.100:12345"

		result := app.GetClientIP(req)
		if result != "192.168.1.100" {
			t.Errorf("expected 192.168.1.100, got %s", result)
		}
	})

	t.Run("trusted proxy extracts X-Forwarded-For", func(t *testing.T) {
		checker, _ := NewTrustedProxyChecker([]string{"10.0.0.1"})
		app.trustedProxies = checker

		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "10.0.0.1:12345"
		req.Header.Set("X-Forwarded-For", "203.0.113.50, 10.0.0.1")

		result := app.GetClientIP(req)
		if result != "203.0.113.50" {
			t.Errorf("expected 203.0.113.50, got %s", result)
		}
	})

	t.Run("untrusted proxy ignores headers", func(t *testing.T) {
		checker, _ := NewTrustedProxyChecker([]string{"10.0.0.1"})
		app.trustedProxies = checker

		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "192.168.1.100:12345"
		req.Header.Set("X-Forwarded-For", "203.0.113.50")

		result := app.GetClientIP(req)
		if result != "192.168.1.100" {
			t.Errorf("expected 192.168.1.100, got %s", result)
		}
	})
}
