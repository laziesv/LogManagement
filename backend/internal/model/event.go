package model

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID        int64           `json:"id"`
	Timestamp time.Time       `json:"@timestamp"`
	Tenant    string          `json:"tenant"`
	Source    string          `json:"source"`
	EventType string          `json:"event_type"`
	Severity  int             `json:"severity"`
	SrcIP     string          `json:"src_ip"`
	User      string          `json:"user"`
	Host      string          `json:"host"`
	Action    string          `json:"action"`
	Fields    map[string]any  `json:"fields"`
	Raw       json.RawMessage `json:"raw"`
}
type Filter struct {
	Tenant, Source, Query string
	From, To              time.Time
	Limit, Offset         int
}
