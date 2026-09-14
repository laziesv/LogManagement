// Package config loads and validates environment configuration.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DatabaseURL    string
	AdminPassword  string
	ViewerPassword string
	APIKeyA        string
	APIKeyB        string
	HTTPAddress    string
	Origin         string
	SecureCookies  bool
	SyslogAddress  string
	SyslogTenant   string
	RetentionDays  int
}

func Load() (Config, error) { return load(os.Getenv) }

func load(getenv func(string) string) (Config, error) {
	value := func(key, fallback string) string {
		if v := getenv(key); v != "" {
			return v
		}
		return fallback
	}
	days, err := strconv.Atoi(value("RETENTION_DAYS", "7"))
	if err != nil || days < 7 {
		return Config{}, fmt.Errorf("RETENTION_DAYS must be at least 7")
	}
	cfg := Config{
		DatabaseURL: getenv("DATABASE_URL"), AdminPassword: getenv("ADMIN_PASSWORD"), ViewerPassword: getenv("VIEWER_PASSWORD"),
		APIKeyA: getenv("API_KEY_A"), APIKeyB: getenv("API_KEY_B"), HTTPAddress: value("HTTP_ADDR", ":3000"),
		Origin: value("APP_ORIGIN", "http://localhost:8080"), SecureCookies: value("COOKIE_SECURE", "false") == "true",
		SyslogAddress: value("SYSLOG_ADDR", ":5514"), SyslogTenant: value("SYSLOG_TENANT", "demo-a"), RetentionDays: days,
	}
	if strings.HasPrefix(cfg.Origin, "https://") && !cfg.SecureCookies {
		return Config{}, fmt.Errorf("HTTPS requires COOKIE_SECURE=true")
	}
	return cfg, nil
}
