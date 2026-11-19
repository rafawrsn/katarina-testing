package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds every environment-driven setting the application needs.
// Nothing here is hardcoded; local dev uses .env, production (Railway)
// relies on real environment variables which always take precedence.
type Config struct {
	AppEnv  string
	AppPort string

	DatabaseURL string

	JWTSecret     string
	JWTExpiryHrs  int

	PaymentGateway      string // e.g. "midtrans"
	PaymentServerKey    string
	PaymentClientKey    string
	PaymentEnvironment  string // "sandbox" | "production"

	ShippingAPIKey       string
	ShippingOriginCity   string
	ShippingOriginProv   string

	CORSAllowedOrigins []string

	RateLimitPerMinute int
}

// Load reads .env (if present) and then environment variables into a
// Config. Missing required values fall back to safe local-dev defaults
// so the server can still boot for local exploration.
func Load() *Config {
	// Ignore the error: in production (Railway) there is no .env file
	// and real env vars are already set.
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", getEnv("PORT", "8090")),

		DatabaseURL: getEnv("DATABASE_URL", "postgresql://postgres:password@localhost:5432/fashion_store"),

		JWTSecret:    getEnv("JWT_SECRET", "change_me"),
		JWTExpiryHrs: getEnvInt("JWT_EXPIRY_HOURS", 72),

		PaymentGateway:     getEnv("PAYMENT_GATEWAY", "midtrans"),
		PaymentServerKey:   getEnv("PAYMENT_SERVER_KEY", ""),
		PaymentClientKey:   getEnv("PAYMENT_CLIENT_KEY", ""),
		PaymentEnvironment: getEnv("PAYMENT_ENVIRONMENT", "sandbox"),

		ShippingAPIKey:     getEnv("SHIPPING_API_KEY", ""),
		ShippingOriginCity: getEnv("SHIPPING_ORIGIN_CITY", "Jakarta Selatan"),
		ShippingOriginProv: getEnv("SHIPPING_ORIGIN_PROVINCE", "DKI Jakarta"),

		CORSAllowedOrigins: splitAndTrim(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:3000,http://127.0.0.1:5173,http://localhost:5174,http://127.0.0.1:3000")),

		RateLimitPerMinute: getEnvInt("RATE_LIMIT_PER_MINUTE", 120),
	}

	return cfg
}

// IsProduction reports whether the app is running in production mode.
func (c *Config) IsProduction() bool {
	return strings.EqualFold(c.AppEnv, "production")
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	n := 0
	for _, ch := range v {
		if ch < '0' || ch > '9' {
			return fallback
		}
		n = n*10 + int(ch-'0')
	}
	return n
}

func splitAndTrim(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
