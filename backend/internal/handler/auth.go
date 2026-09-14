package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"

	"logmanagement/backend/internal/model"
)

func (h *Handler) Login(c fiber.Ctx) error {
	var in model.LoginRequest
	if json.Unmarshal(c.Body(), &in) != nil {
		return fiber.NewError(400, "Invalid JSON")
	}
	u, err := h.store.FindUser(c.Context(), strings.ToLower(strings.TrimSpace(in.Email)))
	hash := []byte(u.Hash)
	if err != nil {
		hash = h.dummyHash
	}
	passErr := bcrypt.CompareHashAndPassword(hash, []byte(in.Password))
	if err != nil || passErr != nil {
		return fiber.NewError(401, "Invalid email or password")
	}
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return err
	}
	token := hex.EncodeToString(b)
	expires := time.Now().Add(8 * time.Hour)
	if err = h.store.CreateSession(c.Context(), token, u.ID, expires); err != nil {
		return err
	}
	c.Cookie(&fiber.Cookie{Name: "session", Value: token, Path: "/", HTTPOnly: true, Secure: h.secureCookies, SameSite: "Strict", Expires: expires})
	return c.JSON(u)
}

func (h *Handler) Me(c fiber.Ctx) error {
	return c.JSON(c.Locals("user"))
}

func (h *Handler) Logout(c fiber.Ctx) error {
	if err := h.store.DeleteSession(c.Context(), c.Cookies("session")); err != nil {
		return err
	}
	c.Cookie(&fiber.Cookie{Name: "session", Value: "", Path: "/", HTTPOnly: true, Secure: h.secureCookies, SameSite: "Strict", Expires: time.Now().Add(-time.Hour), MaxAge: -1})
	return c.SendStatus(204)
}
