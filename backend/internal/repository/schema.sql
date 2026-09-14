CREATE TABLE IF NOT EXISTS tenants (
 id text PRIMARY KEY,
 api_key_hash text UNIQUE NOT NULL,
 alert_enabled boolean NOT NULL DEFAULT true,
 alert_threshold integer NOT NULL DEFAULT 5 CHECK(alert_threshold BETWEEN 2 AND 1000),
 alert_window integer NOT NULL DEFAULT 5 CHECK(alert_window BETWEEN 1 AND 60)
);
CREATE TABLE IF NOT EXISTS users (
 id text PRIMARY KEY,
 tenant text NOT NULL REFERENCES tenants(id),
 email text UNIQUE NOT NULL,
 role text NOT NULL CHECK(role IN ('admin','viewer')),
 password_hash text NOT NULL
);
CREATE TABLE IF NOT EXISTS sessions (
 token_hash text PRIMARY KEY,
 user_id text NOT NULL REFERENCES users(id),
 expires_at timestamptz NOT NULL
);
CREATE TABLE IF NOT EXISTS logs (
 id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
 timestamp timestamptz NOT NULL,
 ingested_at timestamptz NOT NULL DEFAULT now(),
 tenant text NOT NULL REFERENCES tenants(id),
 source text NOT NULL,
 event_type text NOT NULL,
 severity integer NOT NULL CHECK(severity BETWEEN 0 AND 10),
 src_ip text NOT NULL DEFAULT '',
 username text NOT NULL DEFAULT '',
 host text NOT NULL DEFAULT '',
 action text NOT NULL DEFAULT '',
 fields jsonb NOT NULL,
 raw jsonb NOT NULL
);
CREATE INDEX IF NOT EXISTS logs_tenant_time ON logs(tenant,timestamp DESC,id DESC);
CREATE INDEX IF NOT EXISTS logs_source_time ON logs(tenant,source,timestamp DESC);
CREATE INDEX IF NOT EXISTS logs_login_window ON logs(tenant,src_ip,timestamp) WHERE event_type='login_failed';
CREATE INDEX IF NOT EXISTS logs_fields_gin ON logs USING gin(fields);
CREATE INDEX IF NOT EXISTS logs_retention ON logs(ingested_at);
CREATE TABLE IF NOT EXISTS alerts (
 id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
 tenant text NOT NULL REFERENCES tenants(id),
 src_ip text NOT NULL,
 event_count bigint NOT NULL,
 created_at timestamptz NOT NULL DEFAULT now(),
 status text NOT NULL DEFAULT 'open' CHECK(status IN ('open','acknowledged'))
);
CREATE INDEX IF NOT EXISTS alerts_tenant_time ON alerts(tenant,created_at DESC);
