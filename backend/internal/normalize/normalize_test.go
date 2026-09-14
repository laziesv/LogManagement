package normalize

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNormalizeSources(t *testing.T) {
	now := time.Date(2026, 9, 14, 10, 0, 0, 0, time.UTC)
	cases := []struct{ raw, source, event, ip string }{
		{`{"source":"api","event_type":"app_login_failed","ip":"203.0.113.7"}`, "api", "login_failed", "203.0.113.7"},
		{`{"source":"ad","EventID":4625,"IpAddress":"203.0.113.8","TargetUserName":"eve"}`, "ad", "login_failed", "203.0.113.8"},
		{`{"source":"aws","eventName":"CreateUser","sourceIPAddress":"198.51.100.1","eventTime":"2026-09-14T09:00:00Z"}`, "aws", "CreateUser", "198.51.100.1"},
		{`{"source":"m365","Operation":"UserLoggedIn","ClientIP":"198.51.100.2"}`, "m365", "UserLoggedIn", "198.51.100.2"},
		{`{"source":"crowdstrike","event_type":"malware_detected","severity":8}`, "crowdstrike", "malware_detected", ""},
	}
	for _, c := range cases {
		t.Run(c.source, func(t *testing.T) {
			e, err := Normalize(json.RawMessage(c.raw), "demo-a", now)
			if err != nil {
				t.Fatal(err)
			}
			if e.Source != c.source || e.EventType != c.event || e.SrcIP != c.ip || e.Tenant != "demo-a" {
				t.Fatalf("unexpected event %+v", e)
			}
			if string(e.Raw) != c.raw {
				t.Fatal("raw payload changed")
			}
		})
	}
}
func TestRejectInvalidAndCrossTenant(t *testing.T) {
	for _, raw := range []string{`null`, `[]`, `{"tenant":"demo-b"}`, `{"severity":11}`, `{"severity":1.5}`, `{"@timestamp":"yesterday"}`, `{"src_ip":"not-an-ip"}`, `{"source":"unknown"}`, `{"@timestamp":"2099-01-01T00:00:00Z"}`} {
		if _, err := Normalize([]byte(raw), "demo-a", time.Now()); err == nil {
			t.Errorf("accepted %s", raw)
		}
	}
}
func TestSyslogNormalization(t *testing.T) {
	now := time.Date(2026, 9, 14, 10, 0, 0, 0, time.UTC)
	line := `<134>Sep 14 09:59:00 fw01 vendor=demo product=ngfw action=deny src=10.0.1.10 dst=8.8.8.8 spt=5353 dpt=53 proto=udp msg="DNS blocked" tenant=demo-b`
	e, err := NormalizeSyslog(line, "demo-a", now)
	if err != nil {
		t.Fatal(err)
	}
	if e.Source != "firewall" || e.Tenant != "demo-a" || e.Host != "fw01" || e.SrcIP != "10.0.1.10" {
		t.Fatalf("unexpected %+v", e)
	}
	if e.Fields["msg"] != "DNS blocked" {
		t.Fatal("quoted value lost")
	}
	var raw string
	if json.Unmarshal(e.Raw, &raw) != nil || raw != line {
		t.Fatal("raw syslog not preserved")
	}
}
func TestDecodeBatch(t *testing.T) {
	b, e := DecodeBatch([]byte(`{"Records":[{"eventName":"CreateUser"}]}`))
	if e != nil {
		t.Fatal(e)
	}
	event, e := Normalize(b[0], "demo-a", time.Now())
	if e != nil || event.Source != "aws" {
		t.Fatal("AWS envelope failed")
	}
	for _, bad := range []string{`[]`, `null`, `{"Records":[null]}`, `{"Records":null}`, `123`} {
		if _, err := DecodeBatch([]byte(bad)); err == nil {
			t.Errorf("accepted %s", bad)
		}
	}
}
