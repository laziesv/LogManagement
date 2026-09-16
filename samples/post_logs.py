"""Send sample logs using only Python's standard library. Credentials never printed."""
import argparse
import datetime
import json
import os
from pathlib import Path
import urllib.request

def dotenv(root):
    values = {}
    path = root / '.env'
    if not path.exists():
        return values
    for line in path.read_text(encoding='utf-8').splitlines():
        line = line.strip()
        if not line or line.startswith('#') or '=' not in line:
            continue
        name, value = line.split('=', 1)
        values[name.strip()] = value.strip().strip('"').strip("'")
    return values

p = argparse.ArgumentParser()
p.add_argument('--url', default='http://localhost:8080')
p.add_argument('--file', default=str(Path(__file__).with_name('events.json')))
p.add_argument('--tenant', choices=['a', 'b'], default='a', help='Use API_KEY_A or API_KEY_B from .env when LOG_API_KEY is not set')
p.add_argument('--alert', action='store_true', help='Send five fresh failed logins to trigger the default rule')
args = p.parse_args()
root = Path(__file__).resolve().parents[1]
env = dotenv(root)
key = os.environ.get('LOG_API_KEY') or env.get('API_KEY_' + args.tenant.upper())
if not key:
    raise SystemExit('Set LOG_API_KEY or run from a project with API_KEY_A/API_KEY_B in .env')
now = datetime.datetime.now(datetime.timezone.utc).isoformat()
if args.alert:
    data = [{'source': 'api', 'event_type': 'app_login_failed', 'ip': '203.0.113.7', 'user': 'alice', '@timestamp': now} for _ in range(5)]
else:
    data = json.loads(Path(args.file).read_text(encoding='utf-8'))
req = urllib.request.Request(args.url.rstrip('/') + '/ingest', data=json.dumps(data).encode(), headers={'Content-Type': 'application/json', 'X-API-Key': key}, method='POST')
with urllib.request.urlopen(req, timeout=30) as response:
    print(response.read().decode())
