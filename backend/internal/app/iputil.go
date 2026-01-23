package app

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

// TrustedProxyChecker validates whether a remote address is from a trusted proxy.
type TrustedProxyChecker struct {
	cidrs []*net.IPNet
	ips   []net.IP
}

// NewTrustedProxyChecker creates a checker from a list of IP addresses or CIDR ranges.
// Examples: "127.0.0.1", "10.0.0.0/8", "::1", "fd00::/8"
func NewTrustedProxyChecker(trusted []string) (*TrustedProxyChecker, error) {
	checker := &TrustedProxyChecker{}

	for _, entry := range trusted {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		// Try parsing as CIDR first
		if strings.Contains(entry, "/") {
			_, cidr, err := net.ParseCIDR(entry)
			if err != nil {
				return nil, fmt.Errorf("invalid CIDR %q: %w", entry, err)
			}
			checker.cidrs = append(checker.cidrs, cidr)
		} else {
			// Parse as single IP
			ip := net.ParseIP(entry)
			if ip == nil {
				return nil, fmt.Errorf("invalid IP address %q", entry)
			}
			checker.ips = append(checker.ips, ip)
		}
	}

	return checker, nil
}

// IsTrusted checks if the given IP address is in the trusted list.
func (c *TrustedProxyChecker) IsTrusted(ipStr string) bool {
	if c == nil || (len(c.cidrs) == 0 && len(c.ips) == 0) {
		return false
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// Check exact IP matches
	for _, trusted := range c.ips {
		if trusted.Equal(ip) {
			return true
		}
	}

	// Check CIDR ranges
	for _, cidr := range c.cidrs {
		if cidr.Contains(ip) {
			return true
		}
	}

	return false
}

// HasTrustedProxies returns true if any trusted proxies are configured.
func (c *TrustedProxyChecker) HasTrustedProxies() bool {
	return c != nil && (len(c.cidrs) > 0 || len(c.ips) > 0)
}

// GetClientIP extracts the real client IP from a request.
// If the request comes from a trusted proxy, it extracts from forwarding headers.
// Otherwise, it uses RemoteAddr directly.
// All extracted IPs are validated; invalid values fall back to RemoteAddr.
func (app *Application) GetClientIP(r *http.Request) string {
	remoteIP := extractIPFromAddr(r.RemoteAddr)

	// Nil guard: if trustedProxies wasn't initialized, use RemoteAddr
	if app.trustedProxies == nil {
		return remoteIP
	}

	// Only trust forwarding headers if request came from a trusted proxy
	if app.trustedProxies.IsTrusted(remoteIP) {
		// CF-Connecting-IP is set by Cloudflare
		if ip := parseAndValidateIP(r.Header.Get("CF-Connecting-IP")); ip != "" {
			return ip
		}
		// X-Real-IP is typically set by nginx
		if ip := parseAndValidateIP(r.Header.Get("X-Real-IP")); ip != "" {
			return ip
		}
		// X-Forwarded-For may contain multiple IPs: client, proxy1, proxy2
		// The first IP is the original client
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			firstIP := xff
			if idx := strings.Index(xff, ","); idx != -1 {
				firstIP = xff[:idx]
			}
			if ip := parseAndValidateIP(firstIP); ip != "" {
				return ip
			}
		}
	}

	return remoteIP
}

// parseAndValidateIP validates that a string is a valid IP address.
// Returns the canonical string representation, or empty string if invalid.
func parseAndValidateIP(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	ip := net.ParseIP(s)
	if ip == nil {
		return ""
	}
	return ip.String()
}

// extractIPFromAddr extracts just the IP address from an address string.
// Handles both "ip:port" and "[ipv6]:port" formats.
func extractIPFromAddr(addr string) string {
	// Handle IPv6 with brackets: [::1]:8080
	if strings.HasPrefix(addr, "[") {
		if idx := strings.Index(addr, "]:"); idx != -1 {
			return addr[1:idx]
		}
		// Just [::1] without port
		if strings.HasSuffix(addr, "]") {
			return addr[1 : len(addr)-1]
		}
	}

	// Handle IPv4 or IPv6 without brackets
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		// No port, return as-is
		return addr
	}
	return host
}
