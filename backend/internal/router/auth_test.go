package router

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"logmanagement/backend/internal/model"
	"logmanagement/backend/internal/repository"
)

type sessionStore struct {
	repository.Store
	user  model.User
	token string
}

func (s *sessionStore) FindUser(_ context.Context, email string) (model.User, error) {
	if email != s.user.Email {
		return model.User{}, repository.ErrUnauthorized
	}
	return s.user, nil
}
func (s *sessionStore) CreateSession(_ context.Context, token, id string, _ time.Time) error {
	if id != s.user.ID {
		return repository.ErrUnauthorized
	}
	s.token = token
	return nil
}
func (s *sessionStore) Session(_ context.Context, token string) (model.User, error) {
	if token == "" || token != s.token {
		return model.User{}, repository.ErrUnauthorized
	}
	return s.user, nil
}
func (s *sessionStore) DeleteSession(_ context.Context, token string) error {
	if token != s.token {
		return repository.ErrUnauthorized
	}
	s.token = ""
	return nil
}

func TestLoginSessionLifecycle(t *testing.T) {
	password := "test-password-123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	store := &sessionStore{user: model.User{ID: "admin-a", Email: "admin.a@demo.local", Tenant: "demo-a", Role: "admin", Hash: string(hash)}}
	server := NewServer(store, Config{Origin: "https://logs.example.com", SecureCookies: true})
	request := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(`{"email":"admin.a@demo.local","password":"test-password-123"}`))
	request.Header.Set("Origin", "https://logs.example.com")
	request.Header.Set("Content-Type", "application/json")
	response, err := server.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("login %d: %s", response.StatusCode, body)
	}
	var result map[string]any
	if err = json.Unmarshal(body, &result); err != nil {
		t.Fatal(err)
	}
	if result["tenant"] != "demo-a" || result["role"] != "admin" {
		t.Fatalf("wrong user JSON: %s", body)
	}
	if _, ok := result["Hash"]; ok {
		t.Fatal("password hash exposed")
	}
	if strings.Contains(string(body), string(hash)) || strings.Contains(string(body), password) {
		t.Fatal("credentials exposed in response")
	}
	var cookie *http.Cookie
	for _, c := range response.Cookies() {
		if c.Name == "session" {
			cookie = c
		}
	}
	if cookie == nil || cookie.Value == "" || !cookie.HttpOnly || !cookie.Secure || cookie.SameSite != http.SameSiteStrictMode {
		t.Fatal("invalid session cookie")
	}
	for _, step := range []struct {
		method, path string
		status       int
	}{{"GET", "/api/auth/me", 200}, {"POST", "/api/auth/logout", 204}, {"GET", "/api/auth/me", 401}} {
		req := httptest.NewRequest(step.method, step.path, nil)
		req.AddCookie(cookie)
		req.Header.Set("Origin", "https://logs.example.com")
		res, err := server.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != step.status {
			t.Fatalf("%s got %d want %d", step.path, res.StatusCode, step.status)
		}
	}
}
