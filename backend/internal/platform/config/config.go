// Package config carga la configuración desde variables de entorno.
package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	Env                string
	LogLevel           string
	HTTPAddr           string
	CORSAllowedOrigins []string
	DatabaseURL        string
	RedisURL           string
	JWTSecret          string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
}

// Load lee la configuración y falla si falta algún valor obligatorio.
func Load() (Config, error) {
	cfg := Config{
		Env:                getEnv("APP_ENV", "local"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		HTTPAddr:           getEnv("HTTP_ADDR", ":8080"),
		CORSAllowedOrigins: splitList(getEnv("CORS_ALLOWED_ORIGINS", "")),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		RedisURL:           os.Getenv("REDIS_URL"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
	}

	var err error
	if cfg.AccessTokenTTL, err = getDuration("ACCESS_TOKEN_TTL", 15*time.Minute); err != nil {
		return cfg, err
	}
	if cfg.RefreshTokenTTL, err = getDuration("REFRESH_TOKEN_TTL", 30*24*time.Hour); err != nil {
		return cfg, err
	}

	var missing []string
	if cfg.DatabaseURL == "" {
		missing = append(missing, "DATABASE_URL")
	}
	if cfg.RedisURL == "" {
		missing = append(missing, "REDIS_URL")
	}
	if cfg.JWTSecret == "" {
		missing = append(missing, "JWT_SECRET")
	}
	if len(missing) > 0 {
		return cfg, fmt.Errorf("faltan variables de entorno: %s", strings.Join(missing, ", "))
	}
	return cfg, nil
}

func (c Config) IsLocal() bool { return c.Env == "local" }

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) (time.Duration, error) {
	v := os.Getenv(key)
	if v == "" {
		return fallback, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("%s inválido: %w", key, err)
	}
	return d, nil
}

func splitList(v string) []string {
	var out []string
	for _, item := range strings.Split(v, ",") {
		if item = strings.TrimSpace(item); item != "" {
			out = append(out, item)
		}
	}
	return out
}
