package normalize

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"logmanagement/backend/internal/model"
)

// Canonical names win over aliases. Blank/null optional values are absent.
func firstField(m map[string]any, keys ...string) (any, bool) {
	for _, key := range keys {
		v, ok := m[key]
		if !ok || v == nil || v == "" {
			continue
		}
		return v, true
	}
	return nil, false
}

func textField(m map[string]any, keys ...string) (string, error) {
	v, ok := firstField(m, keys...)
	if !ok {
		return "", nil
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string", keys[0])
	}
	return s, nil
}

func integerField(m map[string]any, min, max int, keys ...string) error {
	v, ok := firstField(m, keys...)
	if !ok {
		delete(m, keys[0])
		return nil
	}
	n, err := strconv.Atoi(fmt.Sprint(v))
	if err != nil || n < min || n > max {
		return fmt.Errorf("%s must be an integer from %d to %d", keys[0], min, max)
	}
	m[keys[0]] = n
	return nil
}

// Extra fields stay in JSONB; raw preserves the original payload and aliases.
func normalizeFields(e *model.Event) error {
	m := e.Fields
	m["source"] = e.Source
	for _, keys := range [][]string{
		{"vendor", "Vendor"}, {"product", "Product"},
		{"dst_ip", "dst", "destinationIPAddress", "DestinationIp"},
		{"url", "request_uri", "RequestUri"},
		{"http_method", "method", "HttpMethod"},
	} {
		s, err := textField(m, keys...)
		if err != nil {
			return err
		}
		if s != "" {
			m[keys[0]] = s
		} else {
			delete(m, keys[0])
		}
	}
	if ip, ok := m["dst_ip"].(string); ok {
		parsed := net.ParseIP(ip)
		if parsed == nil {
			return fmt.Errorf("dst_ip must be an IP address")
		}
		m["dst_ip"] = parsed.String()
	}
	for _, spec := range []struct {
		keys     []string
		min, max int
	}{
		{[]string{"src_port", "spt", "sport", "sourcePort", "IpPort"}, 0, 65535},
		{[]string{"dst_port", "dpt", "dport", "destinationPort", "DestinationPort"}, 0, 65535},
		{[]string{"status_code", "http_status", "statusCode"}, 100, 599},
	} {
		if err := integerField(m, spec.min, spec.max, spec.keys...); err != nil {
			return err
		}
	}
	if v, ok := firstField(m, "protocol", "proto", "Protocol"); ok {
		var p string
		switch n := v.(type) {
		case string:
			p = strings.ToLower(strings.TrimSpace(n))
		case float64:
			if n < 0 || n > 255 || n != float64(int(n)) {
				return fmt.Errorf("protocol number must be an integer from 0 to 255")
			}
			p = strconv.Itoa(int(n))
		default:
			return fmt.Errorf("protocol must be a string or IP protocol number")
		}
		switch p {
		case "6":
			p = "tcp"
		case "17":
			p = "udp"
		case "1":
			p = "icmp"
		case "58":
			p = "icmpv6"
		}
		m["protocol"] = p
	}
	if method, ok := m["http_method"].(string); ok {
		m["http_method"] = strings.ToUpper(strings.TrimSpace(method))
	}
	defaults := map[string][2]string{
		"aws": {"aws", "cloudtrail"}, "m365": {"microsoft", "microsoft-365"},
		"ad": {"microsoft", "windows-security"}, "crowdstrike": {"crowdstrike", "falcon"},
	}
	if pair, ok := defaults[e.Source]; ok {
		if value(m, "vendor") == "" {
			m["vendor"] = pair[0]
		}
		if value(m, "product") == "" {
			m["product"] = pair[1]
		}
	}
	if err := normalizeCloud(m); err != nil {
		return err
	}
	outcome, err := textField(m, "outcome")
	if err != nil {
		return err
	}
	if outcome != "" {
		outcome = strings.ToLower(strings.TrimSpace(outcome))
		if outcome != "success" && outcome != "failure" && outcome != "unknown" {
			return fmt.Errorf("outcome must be success, failure or unknown")
		}
	} else {
		outcome = "unknown"
		switch strings.ToLower(value(m, "ResultStatus", "status")) {
		case "success", "succeeded":
			outcome = "success"
		case "failure", "failed", "error":
			outcome = "failure"
		}
		if e.EventType == "login_failed" || value(m, "errorCode") != "" {
			outcome = "failure"
		} else if outcome == "unknown" && (e.EventType == "login_success" || strings.EqualFold(e.EventType, "UserLoggedIn")) {
			outcome = "success"
		}
		if outcome == "unknown" {
			if status, ok := m["status_code"].(int); ok {
				if status >= 400 {
					outcome = "failure"
				} else if status >= 200 {
					outcome = "success"
				}
			}
		}
	}
	e.Action = strings.ToLower(strings.TrimSpace(e.Action))
	if e.Action == "" && (e.EventType == "login_failed" || e.EventType == "login_success" || strings.EqualFold(e.EventType, "UserLoggedIn")) {
		e.Action = "login"
	}
	m["outcome"] = outcome
	// The same canonical values are exposed inside fields and in indexed columns.
	m["@timestamp"] = e.Timestamp.Format(time.RFC3339Nano)
	m["tenant"], m["source"], m["event_type"] = e.Tenant, e.Source, e.EventType
	m["severity"], m["src_ip"], m["user"], m["host"], m["action"] = e.Severity, e.SrcIP, e.User, e.Host, e.Action
	return nil
}

func normalizeCloud(m map[string]any) error {
	cloud := map[string]any{}
	if v, ok := m["cloud"]; ok && v != nil {
		var valid bool
		cloud, valid = v.(map[string]any)
		if !valid {
			return fmt.Errorf("cloud must be an object")
		}
	}
	for _, spec := range []struct{ key, alias string }{
		{"account_id", "recipientAccountId"}, {"region", "awsRegion"}, {"service", "eventSource"},
	} {
		candidates := map[string]any{"value": cloud[spec.key], "dotted": m["cloud."+spec.key]}
		if value(m, "source") == "aws" {
			candidates["alias"] = m[spec.alias]
		}
		s, err := textField(candidates, "value", "dotted", "alias")
		if err != nil {
			return fmt.Errorf("cloud.%s must be a string", spec.key)
		}
		if s != "" {
			if spec.key == "service" && value(m, "source") == "aws" {
				s = strings.TrimSuffix(strings.ToLower(s), ".amazonaws.com")
			}
			cloud[spec.key] = s
		} else {
			delete(cloud, spec.key)
		}
	}
	if len(cloud) > 0 {
		m["cloud"] = cloud
	} else {
		delete(m, "cloud")
	}
	return nil
}
