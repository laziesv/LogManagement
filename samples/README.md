# Two-company samples

| Tenant | Fictional company | Admin | Viewer | Ingest credential |
|---|---|---|---|---|
| demo-a | Atlas Technology | admin.a@demo.local | viewer.a@demo.local | API_KEY_A |
| demo-b | Beacon Retail | admin.b@demo.local | viewer.b@demo.local | API_KEY_B |

Both Admins use ADMIN_PASSWORD from the root .env. Both Viewers use VIEWER_PASSWORD. Each Admin belongs only to its own tenant. The company field is a descriptive label; authorization is enforced by the credential's tenant, not by that label.

Each tenant folder contains:
- events.json: seven events covering API, Firewall, Network, CrowdStrike, AWS, Microsoft 365 and Windows/AD.
- abnormal.json: five failed logins from one IP, sufficient for the default enabled 5-in-5-minute alert rule.

Atlas uses 10.10.x.x, atlas-prefixed hosts, atlas.example users, and cloud account 111111111111. Beacon uses 10.20.x.x, beacon-prefixed hosts, beacon.example users, and cloud account 222222222222. All data is synthetic. Timestamps are omitted so ingestion uses the current time.

Run from the project root with the stack running:

```powershell
python samples/seed_companies.py
```

The script reads the root .env, sends 12 events with each tenant's own API key, then signs in as each Viewer to verify data isolation and write restrictions. It appends data; rerunning creates another batch. Existing logs and credentials are preserved. Repeat failed-login samples within the alert cooldown may reuse the existing alert; changed/disabled rules may not trigger.

For manual file upload, use an Admin of the matching tenant. Use admin.a@demo.local for demo-a and admin.b@demo.local for demo-b. Neither Viewer can upload. Changing only a payload's tenant while using the other tenant's credential is rejected.

Open Log explorer with the matching Viewer, select all sources and the last 24 hours, and search `company-demo`, `Atlas Technology` or `Beacon Retail`. The sidebar still uses the existing tenant ID; this sample setup does not rename tenant records. Generic samples and the real Nginx feed may also appear in demo-a unless filtered.
