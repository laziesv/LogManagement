"""Load two synthetic companies into their own tenants, then verify Viewer isolation."""
import http.cookiejar
import json
from pathlib import Path
import urllib.error
import urllib.parse
import urllib.request
import uuid

root = Path(__file__).resolve().parents[1]
config = dict(line.split('=', 1) for line in (root / '.env').read_text().splitlines() if line and not line.startswith('#'))
base = config.get('APP_ORIGIN', 'http://localhost:8080').rstrip('/')
run_id = 'companies-' + uuid.uuid4().hex

def request(path, data=None, *, key=None, session=None, expected=200, method=None):
    headers = {'Content-Type': 'application/json', 'Origin': base}
    if key:
        headers['X-API-Key'] = key
    req = urllib.request.Request(base + path, data=None if data is None else json.dumps(data).encode(), headers=headers, method=method)
    try:
        response = (session or urllib.request.build_opener()).open(req, timeout=20)
    except urllib.error.HTTPError as error:
        response = error
    with response:
        if response.status != expected:
            raise RuntimeError(f'{path}: HTTP {response.status}, expected {expected}')
        if response.status == 204:
            return None
        return json.load(response)

# Validate both filesets before sending either tenant's batch.
batches = {}
for suffix in ('a', 'b'):
    records = []
    for filename in ('events.json', 'abnormal.json'):
        records.extend(json.loads((root / 'samples' / ('tenant-' + suffix) / filename).read_text(encoding='utf-8')))
    assert len(records) == 12 and all(item['tenant'] == 'demo-' + suffix for item in records)
    for record in records:
        record['sample_run'] = run_id
    batches[suffix] = records

for suffix, records in batches.items():
    result = request('/api/ingest', records, key=config['API_KEY_' + suffix.upper()], expected=201)
    assert result['tenant'] == 'demo-' + suffix and result['accepted'] == len(records)
    print(f"Loaded {len(records)} events: {records[0]['company']} (demo-{suffix})")

for suffix in ('a', 'b'):
    session = urllib.request.build_opener(urllib.request.HTTPCookieProcessor(http.cookiejar.CookieJar()))
    request('/api/auth/login', {'email': f'viewer.{suffix}@demo.local', 'password': config['VIEWER_PASSWORD']}, session=session)
    result = request('/api/logs?' + urllib.parse.urlencode({'q': run_id}), session=session)
    assert result['total'] == 12
    assert all(item['tenant'] == 'demo-' + suffix and item['fields']['company'] == batches[suffix][0]['company'] for item in result['items'])
    other = 'b' if suffix == 'a' else 'a'
    request('/api/logs?tenant=demo-' + other, session=session, expected=403)
    request('/api/ingest', batches[suffix][:1], session=session, expected=403)
    print(f'PASS: Viewer {suffix.upper()} sees only its own company and cannot ingest or access the other tenant.')
admin_b = urllib.request.build_opener(urllib.request.HTTPCookieProcessor(http.cookiejar.CookieJar()))
request('/api/auth/login', {'email': 'admin.b@demo.local', 'password': config['ADMIN_PASSWORD']}, session=admin_b)
rule = request('/api/rule', session=admin_b)
request('/api/rule', rule, session=admin_b, expected=204, method='PUT')
request('/api/logs?tenant=demo-a', session=admin_b, expected=403)
print('PASS: Admin B can sign in, manage demo-b rules, and cannot access demo-a.')
print('Search company-demo or the company name in Log explorer. Existing data is preserved.')
