package app

import "time"

// Config holds all application configuration.
// Configuration is loaded via flags in main.go with environment variable fallbacks.
type Config struct {
	Port int
	Env  string
	DB   struct {
		DSN          string
		MaxOpenConns int
		// MinIdleConns maps to pgxpool.MinConns (minimum idle connections to keep).
		// pgxpool does not support a "max idle" cap.
		MinIdleConns int
		MaxIdleTime  time.Duration
	}
	// TrustedProxies is a list of IP addresses or CIDR ranges that are trusted
	// to set forwarding headers (X-Forwarded-For, X-Real-IP, CF-Connecting-IP).
	// If empty, only RemoteAddr is used (safe default for direct connections).
	TrustedProxies []string
}
