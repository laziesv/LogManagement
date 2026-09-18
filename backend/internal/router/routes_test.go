package router

import (
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"logmanagement/backend/internal/model"
	"logmanagement/backend/internal/repository"
)

type fakeStore struct {
	repository.Store
	role     string
	inserted []model.Event
	filter   model.Filter
}

func (f *fakeStore) Session(_ context.Context, token string) (model.User, error) {
	if token != "valid" {
		return model.User{}, repository.ErrUnauthorized
	}
	return model.User{ID: "u", Tenant: "demo-a", Role: f.role}, nil
}
func (f *fakeStore) APIKeyTenant(_ context.Context, key string) (string, error) {
	if key != "key-a" {
		return "", repository.ErrUnauthorized
	}
	return "demo-a", nil
}
func (f *fakeStore) Insert(_ context.Context, e []model.Event) error {
	f.inserted = append(f.inserted, e...)
	return nil
}
func (f *fakeStore) Logs(_ context.Context, filter model.Filter) ([]model.Event, int64, error) {
	f.filter = filter
	return []model.Event{}, 0, nil
}
func TestHTTPAuthorization(t *testing.T) {
	f := &fakeStore{role: "viewer"}
	s := NewServer(f, Config{Origin: "http://localhost:8080"})
	cases := []struct {
		method, path, body, cookie, key, origin string
		want                                    int
	}{
		{"GET", "/api/logs", "", "", "", "", 401},
		{"GET", "/api/logs?tenant=demo-b", "", "valid", "", "", 403},
		{"GET", "/api/logs?limit=10000", "", "valid", "", "", 400},
		{"GET", "/api/logs", "", "valid", "", "", 200},
		{"POST", "/api/ingest", `{"source":"api"}`, "valid", "", "http://localhost:8080", 403},
		{"PUT", "/api/rule", `{"threshold":5,"window_minutes":5}`, "valid", "", "http://localhost:8080", 403},
		{"POST", "/api/ingest", `{"tenant":"demo-b"}`, "", "key-a", "", 400},
		{"POST", "/api/ingest", `{"source":"api"}`, "", "wrong-key", "", 401},
		{"POST", "/api/ingest", `{"source":"api"}`, "", "key-a", "", 201},
		{"POST", "/ingest", `{"source":"api"}`, "", "key-a", "", 201},
		{"POST", "/api/ingest", `{"source":"api"}`, "valid", "", "https://evil.example", 403},
	}
	for _, c := range cases {
		t.Run(c.method+c.path+c.key+c.origin, func(t *testing.T) {
			req := httptest.NewRequest(c.method, c.path, strings.NewReader(c.body))
			req.Header.Set("Content-Type", "application/json")
			if c.cookie != "" {
				req.Header.Set("Cookie", "session="+c.cookie)
			}
			if c.key != "" {
				req.Header.Set("X-API-Key", c.key)
			}
			if c.origin != "" {
				req.Header.Set("Origin", c.origin)
			}
			res, err := s.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if res.StatusCode != c.want {
				b, _ := io.ReadAll(res.Body)
				t.Fatalf("got %d want %d: %s", res.StatusCode, c.want, b)
			}
		})
	}
	if f.filter.Tenant != "demo-a" {
		t.Fatal("query did not use authenticated tenant")
	}
}
func TestBatchAtomicValidation(t *testing.T) {
	f := &fakeStore{role: "admin"}
	s := NewServer(f, Config{Origin: "http://localhost:8080"})
	req := httptest.NewRequest("POST", "/api/ingest", strings.NewReader(`[{"source":"api"},{"severity":100}]`))
	req.Header.Set("X-API-Key", "key-a")
	res, err := s.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 400 || len(f.inserted) != 0 {
		t.Fatal("invalid batch partially written")
	}
}
func TestFilterTimeRange(t *testing.T) {
	f := &fakeStore{role: "viewer"}
	s := NewServer(f, Config{})
	req := httptest.NewRequest("GET", "/api/logs?from=2026-01-01T00:00:00Z&to=2026-09-01T00:00:00Z", nil)
	req.Header.Set("Cookie", "session=valid")
	res, err := s.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 400 {
		t.Fatal("unbounded range accepted")
	}
}
func TestAdminIngest(t *testing.T) {
	f := &fakeStore{role: "admin"}
	s := NewServer(f, Config{Origin: "http://localhost:8080"})
	body, _ := json.Marshal(map[string]any{"source": "api", "@timestamp": time.Now().UTC().Format(time.RFC3339)})
	foreign := httptest.NewRequest("POST", "/api/ingest", strings.NewReader(string(body)))
	foreign.Header.Set("Cookie", "session=valid")
	foreign.Header.Set("Origin", "https://other.example")
	foreignResponse, err := s.Test(foreign)
	if err != nil {
		t.Fatal(err)
	}
	foreignResponse.Body.Close()
	if foreignResponse.StatusCode != 403 || len(f.inserted) != 0 {
		t.Fatal("cross-origin admin write was accepted")
	}
	req := httptest.NewRequest("POST", "/api/ingest", strings.NewReader(string(body)))
	req.Header.Set("Cookie", "session=valid")
	req.Header.Set("Origin", "http://localhost:8080")
	res, err := s.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 201 || len(f.inserted) != 1 || f.inserted[0].Tenant != "demo-a" {
		t.Fatal("admin ingest failed")
	}
}
