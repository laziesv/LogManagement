"""Check real Nginx -> Syslog -> database flow against a running Compose stack."""
import http.cookiejar
import json
import time
import urllib.error
import urllib.parse
import urllib.request
import uuid
from pathlib import Path

config = dict(line.split('=', 1) for line in (Path(__file__).resolve().parents[1] / '.env').read_text().splitlines() if line and not line.startswith('#'))
base = config.get('APP_ORIGIN', 'http://localhost:8080')

def client(email):
    session = urllib.request.build_opener(urllib.request.HTTPCookieProcessor(http.cookiejar.CookieJar()))
    password = config['ADMIN_PASSWORD'] if email.startswith('admin') else config['VIEWER_PASSWORD']
    req = urllib.request.Request(base + '/api/auth/login', data=json.dumps({'email': email, 'password': password}).encode(), headers={'Content-Type': 'application/json', 'Origin': base})
    with session.open(req, timeout=20) as response:
        assert response.status == 200
    return session

def search(session, query):
    with session.open(base + '/api/logs?' + urllib.parse.urlencode({'q': query, 'source': 'network'}), timeout=20) as response:
        return json.load(response)['items']

admin = client('admin.a@demo.local')
marker = 'real-nginx-' + uuid.uuid4().hex
path = '/assets/' + marker + '.js'
try:
    urllib.request.urlopen(base + path + '?omit_this_query=private-test-marker', timeout=20)
except urllib.error.HTTPError as error:
    assert error.code == 404
    error.close()
else:
    raise AssertionError('Missing static asset should return 404')

deadline = time.monotonic() + 15
items = []
while time.monotonic() < deadline:
    items = search(admin, marker)
    if items:
        break
    time.sleep(0.25)
assert len(items) == 1, 'Expected exactly one real access event'
event = items[0]
assert event['event_type'] == 'nginx_access'
assert event['tenant'] == 'demo-a'
assert event['fields']['status_code'] == 404
assert event['fields']['url'] == path
assert event['fields']['origin'] == 'live'
assert 'private-test-marker' not in json.dumps(event)
assert 'logdesk_access:' in event['raw']
viewer_b = client('viewer.b@demo.local')
assert not search(viewer_b, marker), 'Infrastructure logs leaked to another tenant'
# A request through the API must not feed itself into the collector.
api_marker = 'excluded-' + uuid.uuid4().hex
try:
    admin.open(base + '/api/' + api_marker, timeout=20)
except urllib.error.HTTPError as error:
    error.close()
assert not search(admin, api_marker), 'API access logging should be disabled'
print('PASS: real HTTP 404 -> Nginx -> Syslog -> searchable event; query omitted; tenant isolation verified.')
print('Search nginx_access in Log explorer (source network) to inspect real traffic.')
