"""Integration check for retention cleanup in an isolated temporary schema."""
import subprocess
import uuid

schema = 'retention_test_' + uuid.uuid4().hex


def psql(sql, *, tuples_only=True):
    args = [
        'docker', 'compose', 'exec', '-T', 'postgres',
        'psql', '-U', 'logmanagement', '-d', 'logmanagement',
    ]
    if tuples_only:
        args.extend(['-t', '-A'])
    args.extend(['-c', sql])
    result = subprocess.run(args, check=True, capture_output=True, text=True)
    return result.stdout.strip()


try:
    psql(f'create schema {schema};', tuples_only=False)
    psql(f"""
    create table {schema}.logs (
      id bigint generated always as identity primary key,
      timestamp timestamptz not null,
      ingested_at timestamptz not null default now(),
      tenant text not null,
      source text not null,
      event_type text not null,
      severity integer not null check(severity between 0 and 10),
      src_ip text not null default '',
      username text not null default '',
      host text not null default '',
      action text not null default '',
      fields jsonb not null,
      raw jsonb not null
    );
    """, tuples_only=False)
    psql(f"""
    insert into {schema}.logs(timestamp, ingested_at, tenant, source, event_type, severity, fields, raw)
    values
      (now(), now() - interval '8 days', 'demo-a', 'api', 'retention_old', 3, '{{}}', '{{}}'),
      (now(), now() - interval '6 days', 'demo-a', 'api', 'retention_fresh', 3, '{{}}', '{{}}');
    """, tuples_only=False)
    psql(f"delete from {schema}.logs where ingested_at < now()-(7 * interval '1 day');", tuples_only=False)
    remaining = psql(f"select event_type from {schema}.logs order by event_type;")
    assert remaining == 'retention_fresh', remaining
    print('PASS: retention cleanup deletes logs ingested before 7 days and keeps newer logs.')
finally:
    if schema.startswith('retention_test_'):
        subprocess.run([
            'docker', 'compose', 'exec', '-T', 'postgres',
            'psql', '-U', 'logmanagement', '-d', 'logmanagement',
            '-c', f'drop schema if exists {schema} cascade;',
        ], check=False, capture_output=True, text=True)
