"""Integration smoke tests against the running stack. Read root .env locally."""
import datetime
import http.cookiejar
import json
from pathlib import Path
import urllib.error
import urllib.request

config = dict(line.split('=', 1) for line in (Path(__file__).resolve().parents[1] / '.env').read_text().splitlines() if line and not line.startswith('#'))
base = config.get('APP_ORIGIN', 'http://localhost:8080')

def request(path, data=None, *, opener=None, key=None, expected=200, method=None):
    headers = {'Content-Type': 'application/json', 'Origin': base}
    if key: headers['X-API-Key'] = key
    req = urllib.request.Request(base + path, data=None if data is None else json.dumps(data).encode(), headers=headers, method=method)
    try:
        response = (opener or urllib.request.build_opener()).open(req, timeout=30)
    except urllib.error.HTTPError as e:
        response = e
    with response:
        body = response.read()
        assert response.code == expected, (path, response.code, expected, body)
        return json.loads(body) if body and response.code != 204 else None

def login(email, password):
    opener = urllib.request.build_opener(urllib.request.HTTPCookieProcessor(http.cookiejar.CookieJar()))
    request('/api/auth/login', {'email': email, 'password': password}, opener=opener)
    return opener

request('/api/health')
request('/api/logs', expected=401)
admin = login('admin@demo.local', config['ADMIN_PASSWORD'])
viewer = login('viewer.a@demo.local', config['VIEWER_PASSWORD'])
viewer_b = login('viewer.b@demo.local', config['VIEWER_PASSWORD'])
stamp = datetime.datetime.now(datetime.timezone.utc).isoformat()
marker = 'smoke_' + stamp
request('/api/ingest', {'source':'api','event_type':marker}, key=config['API_KEY_A'], expected=201)
request('/api/ingest', {'source':'api','event_type':'tenant_b_' + marker}, key=config['API_KEY_B'], expected=201)
request('/api/ingest', {'tenant':'demo-b'}, key=config['API_KEY_A'], expected=400)
request('/api/ingest', {'source':'api'}, opener=viewer, expected=403)
request('/api/logs?tenant=demo-b', opener=viewer, expected=403)
result = request('/api/logs', opener=viewer)
assert all(e['tenant'] == 'demo-a' for e in result['items'])
assert any(e['event_type'] == marker for e in result['items'])
result = request('/api/logs', opener=viewer_b)
assert all(e['tenant'] == 'demo-b' for e in result['items'])
assert any(e['event_type'] == 'tenant_b_' + marker for e in result['items'])
old_rule = request('/api/rule', opener=admin)
try:
    request('/api/rule', {'enabled':True,'threshold':5,'window_minutes':5}, opener=admin, method='PUT')
    before = request('/api/alerts', opener=admin)
    occupied = {a['src_ip'] for a in before}
    ip = next(f'192.0.2.{i}' for i in range(1, 255) if f'192.0.2.{i}' not in occupied)
    request('/api/ingest', [{'source':'api','event_type':'login_failed','src_ip':ip,'@timestamp':stamp} for _ in range(5)], key=config['API_KEY_A'], expected=201)
    alerts = request('/api/alerts', opener=admin)
    alert = next(a for a in alerts if a['src_ip'] == ip and a['status'] == 'open')
    request(f"/api/alerts/{alert['id']}/acknowledge", {}, opener=viewer, expected=403)
    request(f"/api/alerts/{alert['id']}/acknowledge", {}, opener=admin, expected=204)
finally:
    request('/api/rule', old_rule, opener=admin, method='PUT')
request('/api/auth/logout', {}, opener=admin, expected=204)
request('/api/auth/me', opener=admin, expected=401)
print('PASS: health, auth, ingestion, tenant isolation, Viewer permissions, alert trigger/acknowledgement, logout')
