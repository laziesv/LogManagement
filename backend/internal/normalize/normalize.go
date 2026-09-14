package normalize

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"logmanagement/backend/internal/model"
)

func value(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil {
			s := fmt.Sprint(v)
			if s != "" {
				return s
			}
		}
	}
	return ""
}
func Normalize(raw json.RawMessage, tenant string, now time.Time) (model.Event, error) {
	var m map[string]any
	if json.Unmarshal(raw, &m) != nil || m == nil {
		return model.Event{}, errors.New("each log must be a JSON object")
	}
	if t := value(m, "tenant"); t != "" && t != tenant {
		return model.Event{}, errors.New("log tenant does not match credential tenant")
	}
	source := strings.ToLower(value(m, "source"))
	if source == "" {
		source = "api"
	}
	switch source {
	case "api", "firewall", "network", "crowdstrike", "aws", "m365", "ad":
	default:
		return model.Event{}, errors.New("unsupported source")
	}
	ts := now.UTC()
	if s := value(m, "@timestamp", "timestamp", "eventTime", "CreationTime"); s != "" {
		var err error
		ts, err = time.Parse(time.RFC3339Nano, s)
		if err != nil {
			return model.Event{}, errors.New("timestamp must be RFC3339")
		}
	}
	if ts.After(now.Add(5 * time.Minute)) {
		return model.Event{}, errors.New("timestamp cannot exceed server time by more than five minutes")
	}
	typ := value(m, "event_type", "eventName", "Operation", "event")
	if typ == "" {
		typ = "log"
	}
	if value(m, "event_id", "EventID") == "4625" || strings.EqualFold(typ, "LogonFailed") || strings.EqualFold(typ, "app_login_failed") {
		typ = "login_failed"
	}
	if value(m, "event_id", "EventID") == "4624" {
		typ = "login_success"
	}
	severity := 3
	if s := value(m, "severity"); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil || n < 0 || n > 10 {
			return model.Event{}, errors.New("severity must be an integer from 0 to 10")
		}
		severity = n
	}
	if typ == "login_failed" && value(m, "severity") == "" {
		severity = 5
	}
	ip := value(m, "src_ip", "ip", "src", "sourceIPAddress", "ClientIP", "IpAddress")
	if ip != "" && net.ParseIP(ip) == nil {
		return model.Event{}, errors.New("src_ip must be an IP address")
	}
	e := model.Event{Timestamp: ts.UTC(), Tenant: tenant, Source: source, EventType: typ, Severity: severity, SrcIP: ip, User: value(m, "user", "UserId", "TargetUserName"), Host: value(m, "host", "Computer"), Action: value(m, "action"), Fields: m, Raw: append(json.RawMessage(nil), raw...)}
	if identity, ok := m["userIdentity"].(map[string]any); ok && e.User == "" {
		e.User = value(identity, "userName", "arn")
	}
	if source == "aws" {
		if _, exists := m["cloud"]; !exists {
			m["cloud"] = map[string]any{"region": m["awsRegion"], "service": m["eventSource"], "account_id": m["recipientAccountId"]}
		}
	}
	return e, nil
}

var kvPattern = regexp.MustCompile(`([A-Za-z_][A-Za-z0-9_]*)=("[^"]*"|\S+)`)
var priPattern = regexp.MustCompile(`^<(\d{1,3})>`)
var rfc3164Pattern = regexp.MustCompile(`^([A-Z][a-z]{2}\s+\d{1,2}\s+\d{2}:\d{2}:\d{2})\s+(\S+)\s+`)

// Syslog uses a fixed deployment tenant, never an untrusted tenant field.
// RFC3164 has no year/timezone; assume UTC and the nearest nonfuture year.
func NormalizeSyslog(line, tenant string, now time.Time) (model.Event, error) {
	m := map[string]any{"source": "network"}
	body := strings.TrimSpace(line)
	if match := priPattern.FindStringSubmatch(body); len(match) > 0 {
		pri, _ := strconv.Atoi(match[1])
		if pri > 191 {
			return model.Event{}, errors.New("invalid syslog priority")
		}
		m["severity"] = (7 - pri%8) * 10 / 7
		body = body[len(match[0]):]
	}
	if match := rfc3164Pattern.FindStringSubmatch(body); len(match) > 0 {
		stamp, err := time.Parse("2006 Jan 2 15:04:05", fmt.Sprintf("%d %s", now.Year(), strings.Join(strings.Fields(match[1]), " ")))
		if err == nil {
			if stamp.After(now.Add(24 * time.Hour)) {
				stamp = stamp.AddDate(-1, 0, 0)
			}
			m["@timestamp"] = stamp.Format(time.RFC3339)
		}
		m["host"] = match[2]
	}
	for _, match := range kvPattern.FindAllStringSubmatch(body, -1) {
		if match[1] != "tenant" {
			m[match[1]] = strings.Trim(match[2], `"`)
		}
	}
	if value(m, "product") == "ngfw" || value(m, "action") == "deny" || value(m, "action") == "allow" {
		m["source"] = "firewall"
	}
	if value(m, "event_type", "event") == "" {
		m["event_type"] = "syslog"
	}
	raw, _ := json.Marshal(m)
	e, err := Normalize(raw, tenant, now)
	e.Raw, _ = json.Marshal(line)
	return e, err
}
