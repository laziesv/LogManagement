package model

import (
	"time"
)

type Alert struct {
	ID        int64     `json:"id"`
	Tenant    string    `json:"tenant"`
	IP        string    `json:"src_ip"`
	Count     int64     `json:"count"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
}
type Rule struct {
	Enabled       bool `json:"enabled"`
	Threshold     int  `json:"threshold"`
	WindowMinutes int  `json:"window_minutes"`
}
