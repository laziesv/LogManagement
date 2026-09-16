"""Integration checks for expanded normalization against the running demo stack."""
import http.cookiejar
import json
import urllib.error
import urllib.parse
import urllib.request
import uuid
from pathlib import Path

root = Path(__file__).resolve().parents[1]
cfg = dict(line.split('=', 1) for line in (root / '.env').read_text().splitlines() if line and not line.startswith('#'))
base = cfg['APP_ORIGIN']
client = urllib.request.build_opener(urllib.request.HTTPCookieProcessor(http.cookiejar.CookieJar()))
def request(path, data=None, expected=200):
    req = urllib.request.Request(base + path, data=json.dumps(data).encode() if data is not None else None, headers={'Content-Type': 'application/json', 'Origin': base})
    try:
        response = client.open(req, timeout=20)
    except urllib.error.HTTPError as error:
        response = error
    with response:
        result = json.load(response)
        assert response.code == expected, (response.code, expected)
        return result

def search(q):
    return request('/api/logs?' + urllib.parse.urlencode({'q': q}))['items']

request('/api/auth/login', {'email': 'admin.a@demo.local', 'password': cfg['ADMIN_PASSWORD']})
marker = 'canonical-' + uuid.uuid4().hex
records = [
    {'source':'firewall','vendor':'demo','product':'ngfw','event_type':'connection_blocked','src':'10.0.1.10','dst':'8.8.8.8','spt':'5353','dpt':'53','proto':'UDP','action':'deny'},
    {'source':'api','event_type':'http_request','method':'get','request_uri':'/example','http_status':'404','protocol':'HTTP/1.1'},
    {'source':'aws','eventName':'CreateUser','eventSource':'iam.amazonaws.com','recipientAccountId':'001234567890','awsRegion':'ap-southeast-1','errorCode':'AccessDenied'},
    {'source':'ad','EventID':4625,'TargetUserName':'demo-user','IpAddress':'192.0.2.88','IpPort':'53001','Computer':'DC01'},
]
for record in records:
    record['test_marker'] = marker
request('/api/ingest', records, expected=201)
items = search(marker)
assert len(items) == 4
by_source = {item['source']: item for item in items}
firewall = by_source['firewall']['fields']
assert firewall['src_port'] == 5353 and firewall['dst_port'] == 53 and firewall['protocol'] == 'udp'
assert by_source['api']['fields']['status_code'] == 404
assert by_source['api']['fields']['outcome'] == 'failure'
assert by_source['aws']['fields']['cloud']['account_id'] == '001234567890'
assert by_source['aws']['fields']['cloud']['service'] == 'iam'
assert by_source['ad']['action'] == 'login'
assert by_source['ad']['fields']['outcome'] == 'failure'
for record in records:
    assert by_source[record['source']]['raw'] == record
assert any(item['id'] == by_source['aws']['id'] for item in search('cloudtrail')), 'Search missed a derived field absent from raw'
bad_marker = 'reject-' + uuid.uuid4().hex
request('/api/ingest', [{'test_marker': bad_marker}, {'dst_port': 70000}], expected=400)
assert not search(bad_marker), 'Invalid batch partially inserted'
print('PASS: typed fields, provider mappings, raw preservation, derived-field search and atomic rejection.')
