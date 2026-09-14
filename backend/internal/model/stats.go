package model

type Count struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}
type Stats struct {
	Total     int64   `json:"total"`
	Critical  int64   `json:"critical"`
	Sources   int64   `json:"sources"`
	Alerts    int64   `json:"alerts"`
	Timeline  []Count `json:"timeline"`
	TopIPs    []Count `json:"top_ips"`
	TopUsers  []Count `json:"top_users"`
	TopEvents []Count `json:"top_events"`
}
