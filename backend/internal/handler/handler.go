package handler

import (
	"golang.org/x/crypto/bcrypt"

	"logmanagement/backend/internal/repository"
)

type Handler struct {
	store         repository.Store
	secureCookies bool
	dummyHash     []byte
}

func New(store repository.Store, secureCookies bool) *Handler {
	// A fixed-cost comparison avoids account-existence timing differences.
	dummy, _ := bcrypt.GenerateFromPassword([]byte("invalid-login-password"), bcrypt.DefaultCost)
	return &Handler{store: store, secureCookies: secureCookies, dummyHash: dummy}
}
