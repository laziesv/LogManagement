package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"logmanagement/backend/internal/model"
)

var ErrUnauthorized = errors.New("unauthorized")

func tokenHash(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }

type Store interface {
	FindUser(context.Context, string) (model.User, error)
	CreateSession(context.Context, string, string, time.Time) error
	Session(context.Context, string) (model.User, error)
	DeleteSession(context.Context, string) error
	APIKeyTenant(context.Context, string) (string, error)
	Insert(context.Context, []model.Event) error
	Logs(context.Context, model.Filter) ([]model.Event, int64, error)
	Stats(context.Context, model.Filter) (model.Stats, error)
	Alerts(context.Context, string) ([]model.Alert, error)
	Acknowledge(context.Context, string, int64) (bool, error)
	GetRule(context.Context, string) (model.Rule, error)
	SetRule(context.Context, string, model.Rule) error
	Ping(context.Context) error
}
