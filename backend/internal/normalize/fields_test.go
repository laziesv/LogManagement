package normalize

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestCanonicalFields(t *testing.T) {
	cases := []struct {
		name, raw string
		want      map[string]any
	}{
		{"firewall", `{"source":"firewall","src":"10.0.0.1","dst":"2001:0db8::1","spt":"5353","dpt":53,"proto":"UDP","vendor":"demo","product":"ngfw","action":"DENY"}`, map[string]any{"src_ip": "10.0.0.1", "dst_ip": "2001:db8::1", "src_port": 5353, "dst_port": 53, "protocol": "udp", "action": "deny", "outcome": "unknown"}},
		{"http", `{"source":"network","method":"get","request_uri":"/hello","http_status":"404","protocol":"HTTP/1.1"}`, map[string]any{"http_method": "GET", "url": "/hello", "status_code": 404, "protocol": "http/1.1", "outcome": "failure"}},
		{"http success", `{"status_code":204}`, map[string]any{"outcome": "success"}},
		{"http informational", `{"status_code":101}`, map[string]any{"outcome": "unknown"}},
		{"windows failed", `{"source":"ad","EventID":4625,"IpPort":"0"}`, map[string]any{"vendor": "microsoft", "product": "windows-security", "src_port": 0, "event_type": "login_failed", "action": "login", "outcome": "failure"}},
		{"windows success", `{"source":"ad","EventID":4624}`, map[string]any{"event_type": "login_success", "action": "login", "outcome": "success"}},
		{"m365", `{"source":"m365","Operation":"UserLoggedIn","ResultStatus":"Failed"}`, map[string]any{"vendor": "microsoft", "outcome": "failure"}},
		{"cloud merge", `{"source":"AWS","cloud":{"region":"custom","extra":"keep"},"recipientAccountId":"001234567890","eventSource":"iam.amazonaws.com","awsRegion":"ignored","errorCode":"AccessDenied"}`, map[string]any{"cloud": map[string]any{"region": "custom", "extra": "keep", "account_id": "001234567890", "service": "iam"}, "outcome": "failure", "vendor": "aws"}},
		{"dotted cloud", `{"cloud.account_id":"0001","cloud.region":"test"}`, map[string]any{"cloud": map[string]any{"account_id": "0001", "region": "test"}}},
		{"precedence", `{"dst_port":443,"dpt":80,"protocol":6,"status_code":500,"outcome":"SUCCESS"}`, map[string]any{"dst_port": 443, "protocol": "tcp", "outcome": "success"}},
		{"explicit vendor", `{"source":"crowdstrike","vendor":"custom","product":"custom-agent"}`, map[string]any{"vendor": "custom", "product": "custom-agent", "outcome": "unknown"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			e, err := Normalize(json.RawMessage(c.raw), "demo-a", time.Now())
			if err != nil {
				t.Fatal(err)
			}
			for key, want := range c.want {
				if !reflect.DeepEqual(e.Fields[key], want) {
					t.Errorf("%s = %#v; want %#v", key, e.Fields[key], want)
				}
			}
			if string(e.Raw) != c.raw {
				t.Fatal("raw changed")
			}
			if e.Fields["tenant"] != e.Tenant || e.Fields["action"] != e.Action || e.Fields["source"] != e.Source {
				t.Fatal("canonical fields disagree with indexed columns")
			}
		})
	}
}

func TestInvalidCanonicalFields(t *testing.T) {
	for _, raw := range []string{
		`{"dst_ip":"bad"}`, `{"dst_ip":{}}`, `{"src_port":65536}`, `{"src_port":true}`, `{"dst_port":-1}`, `{"dst_port":2.5}`, `{"dst_port":[]}`,
		`{"status_code":99}`, `{"status_code":600}`, `{"status_code":"Success"}`, `{"http_method":{}}`, `{"protocol":true}`, `{"protocol":6.5}`,
		`{"outcome":"maybe"}`, `{"cloud":[]}`, `{"cloud":{"region":42}}`, `{"cloud.account_id":123}`, `{"vendor":[]}`,
	} {
		if _, err := Normalize(json.RawMessage(raw), "demo-a", time.Now()); err == nil {
			t.Errorf("accepted invalid input: %s", raw)
		}
	}
}

func TestCanonicalSyslogNetworkFields(t *testing.T) {
	e, err := NormalizeSyslog(`<134>fw01 product=ngfw action=deny src=10.0.1.10 dst=8.8.8.8 spt=5353 dpt=53 proto=udp`, "demo-a", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if e.Fields["src_port"] != 5353 || e.Fields["dst_port"] != 53 || e.Fields["dst_ip"] != "8.8.8.8" || e.Fields["protocol"] != "udp" {
		t.Fatalf("bad syslog fields: %+v", e.Fields)
	}
}
