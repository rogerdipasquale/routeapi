package config

import (
	"log/slog"
	"os"
	"strings"
)

// Config holds runtime configuration from environment variables (parity with Nest).
type Config struct {
	Port        string
	Environment string
	// AllowedOrigins lists permitted CORS origins; empty means no CORS headers (strict).
	AllowedOrigins []string
	LogLevel       slog.Value
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	origins := strings.TrimSpace(os.Getenv("ALLOWED_ORIGINS"))
	var allowed []string
	if origins != "" {
		for _, o := range strings.Split(origins, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				allowed = append(allowed, o)
			}
		}
	}
	env := os.Getenv("NODE_ENV")
	if env == "" {
		env = os.Getenv("GO_ENV")
	}
	logLevelStr := os.Getenv("LOG_LEVEL")
	if logLevelStr == "" {
		logLevelStr = "INFO"
	}
	logLevel := slog.StringValue(logLevelStr)

	return Config{
		Port:           port,
		Environment:    env,
		AllowedOrigins: allowed,
		LogLevel:       logLevel,
	}
}
