"""Insert three related 8-day-old logs into the real logs table for retention demos."""
import subprocess
import uuid

incident_id = 'retention_incident_' + uuid.uuid4().hex
events = [
    ('login_failed', 5, '203.0.113.88', 'alice@example.local', 'auth-01', 'login', 'Initial failed login from test IP'),
    ('password_spray_detected', 8, '203.0.113.88', 'multiple-users', 'auth-01', 'detect', 'Correlated repeated failures for the same source IP'),
    ('session_blocked', 6, '203.0.113.88', 'alice@example.local', 'gateway-01', 'block', 'Gateway blocked the related session attempt'),
]
values = []
for sequence, (event_type, severity, src_ip, user, host, action, message) in enumerate(events, start=1):
    values.append(f"""
  (
    now(), now() - interval '8 days', 'demo-a', 'api', '{incident_id}_{event_type}', {severity},
    '{src_ip}', '{user}', '{host}', '{action}',
    jsonb_build_object(
      'origin','manual-db',
      'retention_case','8_days_old',
      'incident_id','{incident_id}',
      'sequence',{sequence},
      'msg','{message}'
    ),
    jsonb_build_object(
      'source','api',
      'event_type','{incident_id}_{event_type}',
      'retention_case','8_days_old',
      'incident_id','{incident_id}',
      'sequence',{sequence},
      'related_src_ip','{src_ip}',
      'msg','{message}'
    )
  )""")

sql = f"""
with inserted as (
  insert into logs(timestamp, ingested_at, tenant, source, event_type, severity, src_ip, username, host, action, fields, raw)
  values {','.join(values)}
  returning id, tenant, source, event_type, severity, src_ip, username, host, action, timestamp, ingested_at, fields, raw
)
select
  id,
  tenant,
  source,
  event_type,
  severity,
  src_ip,
  username,
  host,
  action,
  fields->>'incident_id' as incident_id,
  fields->>'sequence' as sequence,
  raw
from inserted
order by (fields->>'sequence')::int;

select count(*) as rows_for_incident
from logs
where fields->>'incident_id' = '{incident_id}';
"""

subprocess.run(
    [
        'docker', 'compose', 'exec', '-T', 'postgres',
        'psql', '-U', 'logmanagement', '-d', 'logmanagement',
        '-c', sql,
    ],
    check=True,
)

print()
print('Inserted three related 8-day-old retention test logs.')
print('Search incident before cleanup:')
print(incident_id)
print()
print('With the backend running, retention cleanup should remove it within about 1 minute.')
