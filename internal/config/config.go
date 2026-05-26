// Package config loads ExchangeOS configuration from env vars (12-factor) with optional .env.
// Secrets MUST come from Vault SPI in production — env vars only for local/dev/CI.
package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Env     string
	HTTP    HTTPConfig
	GRPC    GRPCConfig
	DB      DBConfig
	OTel    OTelConfig
	Repos   ReposConfig
}

// ReposConfig selects the repository backend (memory bootstrap vs postgres).
type ReposConfig struct {
	Backend string // "memory" (default) | "postgres"
}

type HTTPConfig struct {
	Port int
}

type GRPCConfig struct {
	Port         int
	MaxRecvBytes int
	MaxSendBytes int
	Reflection   bool
}

type DBConfig struct {
	DSN     string
	MaxConn int
	MinConn int
}

type OTelConfig struct {
	Endpoint   string // OTLP gRPC endpoint
	Sampling   float64
	Insecure   bool
}

// Load reads configuration. Order: process env > .env (if exists) > defaults.
func Load() (*Config, error) {
	// Best-effort .env load (no error if absent).
	_ = godotenv.Load()

	cfg := &Config{
		Env: getEnv("EXCHANGEOS_ENV", "dev"),
		HTTP: HTTPConfig{
			Port: getEnvInt("EXCHANGEOS_HTTP_PORT", 8094),
		},
		GRPC: GRPCConfig{
			Port:         getEnvInt("EXCHANGEOS_GRPC_PORT", 9094),
			MaxRecvBytes: getEnvInt("EXCHANGEOS_GRPC_MAX_RECV", 16<<20), // 16 MiB
			MaxSendBytes: getEnvInt("EXCHANGEOS_GRPC_MAX_SEND", 16<<20),
			Reflection:   getEnvBool("EXCHANGEOS_GRPC_REFLECTION", true),
		},
		DB: DBConfig{
			DSN:     getEnv("EXCHANGEOS_DB_DSN", ""),
			MaxConn: getEnvInt("EXCHANGEOS_DB_MAX_CONN", 25),
			MinConn: getEnvInt("EXCHANGEOS_DB_MIN_CONN", 2),
		},
		OTel: OTelConfig{
			Endpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
			Sampling: getEnvFloat("OTEL_TRACES_SAMPLER_ARG", 1.0),
			Insecure: getEnvBool("OTEL_EXPORTER_OTLP_INSECURE", true),
		},
		Repos: ReposConfig{
			Backend: strings.ToLower(getEnv("EXCHANGEOS_REPO_BACKEND", "memory")),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	if c.HTTP.Port <= 0 || c.HTTP.Port > 65535 {
		return fmt.Errorf("HTTP.Port out of range: %d", c.HTTP.Port)
	}
	if c.GRPC.Port <= 0 || c.GRPC.Port > 65535 {
		return fmt.Errorf("GRPC.Port out of range: %d", c.GRPC.Port)
	}
	if c.HTTP.Port == c.GRPC.Port {
		return fmt.Errorf("HTTP and gRPC ports must differ (both %d)", c.HTTP.Port)
	}
	if c.Env == "production" && c.DB.DSN == "" {
		return fmt.Errorf("EXCHANGEOS_DB_DSN required in production")
	}
	switch c.Repos.Backend {
	case "memory", "postgres":
	default:
		return fmt.Errorf("EXCHANGEOS_REPO_BACKEND must be 'memory' or 'postgres', got %q", c.Repos.Backend)
	}
	if c.Repos.Backend == "postgres" && c.DB.DSN == "" {
		return fmt.Errorf("EXCHANGEOS_DB_DSN required when EXCHANGEOS_REPO_BACKEND=postgres")
	}
	return nil
}

// Redacted returns the DB DSN with the password masked. Use in logs.
func (d DBConfig) Redacted() string {
	if d.DSN == "" {
		return "(empty)"
	}
	u, err := url.Parse(d.DSN)
	if err != nil {
		return "(unparseable)"
	}
	if u.User != nil {
		u.User = url.UserPassword(u.User.Username(), "REDACTED")
	}
	return u.String()
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getEnvFloat(key string, def float64) float64 {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return def
}

func getEnvBool(key string, def bool) bool {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		switch strings.ToLower(v) {
		case "1", "true", "yes", "on":
			return true
		case "0", "false", "no", "off":
			return false
		}
	}
	return def
}
