"""Send sample logs using only Python's standard library. Credentials never printed."""
import argparse
import datetime
import json
import os
from pathlib import Path
import urllib.request

p = argparse.ArgumentParser()
p.add_argument('--url', default='http://localhost:8080')
p.add_argument('--file', default=str(Path(__file__).with_name('events.json')))
p.add_argument('--alert', action='store_true', help='Send five fresh failed logins to trigger the default rule')
args = p.parse_args()
key = os.environ.get('LOG_API_KEY')
if not key:
    raise SystemExit('Set LOG_API_KEY to API_KEY_A or API_KEY_B from .env')
now = datetime.datetime.now(datetime.timezone.utc).isoformat()
if args.alert:
    data = [{'source': 'api', 'event_type': 'app_login_failed', 'ip': '203.0.113.7', 'user': 'alice', '@timestamp': now} for _ in range(5)]
else:
    data = json.loads(Path(args.file).read_text(encoding='utf-8'))
req = urllib.request.Request(args.url.rstrip('/') + '/ingest', data=json.dumps(data).encode(), headers={'Content-Type': 'application/json', 'X-API-Key': key}, method='POST')
with urllib.request.urlopen(req, timeout=30) as response:
    print(response.read().decode())
