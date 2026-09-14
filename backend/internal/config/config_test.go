package config

import "testing"

func TestDefaults(t *testing.T) {
	cfg, err := load(func(string) string { return "" })
	if err != nil {
		t.Fatal(err)
	}
	if cfg.RetentionDays != 7 || cfg.HTTPAddress != ":3000" || cfg.SyslogAddress != ":5514" || cfg.SyslogTenant != "demo-a" || cfg.Origin != "http://localhost:8080" || cfg.SecureCookies {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}

func TestOverrides(t *testing.T) {
	values := map[string]string{"DATABASE_URL": "postgres://example", "ADMIN_PASSWORD": "admin", "VIEWER_PASSWORD": "viewer", "API_KEY_A": "a", "API_KEY_B": "b", "HTTP_ADDR": ":9000", "SYSLOG_ADDR": ":9001", "SYSLOG_TENANT": "tenant-b", "APP_ORIGIN": "https://logs.example.com", "COOKIE_SECURE": "true", "RETENTION_DAYS": "30"}
	cfg, err := load(func(key string) string { return values[key] })
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DatabaseURL != values["DATABASE_URL"] || cfg.AdminPassword != "admin" || cfg.ViewerPassword != "viewer" || cfg.APIKeyA != "a" || cfg.APIKeyB != "b" || cfg.HTTPAddress != ":9000" || cfg.SyslogAddress != ":9001" || cfg.SyslogTenant != "tenant-b" || cfg.Origin != values["APP_ORIGIN"] || cfg.RetentionDays != 30 || !cfg.SecureCookies {
		t.Fatal("environment overrides not preserved")
	}
}

func TestInvalidConfiguration(t *testing.T) {
	for _, values := range []map[string]string{{"RETENTION_DAYS": "6"}, {"RETENTION_DAYS": "bad"}, {"APP_ORIGIN": "https://logs.example.com"}} {
		if _, err := load(func(key string) string { return values[key] }); err == nil {
			t.Fatalf("accepted invalid config %v", values)
		}
	}
}
