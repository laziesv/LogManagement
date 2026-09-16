package repository

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"logmanagement/backend/internal/model"
)

//go:embed schema.sql
var schema string

type Postgres struct{ Pool *pgxpool.Pool }

func OpenPostgres(ctx context.Context, url string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	if _, err = pool.Exec(ctx, schema); err != nil {
		pool.Close()
		return nil, err
	}
	return &Postgres{pool}, nil
}
func (p *Postgres) Seed(ctx context.Context, adminPass, viewerPass, keyA, keyB string) error {
	if len(adminPass) < 12 || len(viewerPass) < 12 || len(adminPass) > 72 || len(viewerPass) > 72 || len(keyA) < 24 || len(keyB) < 24 || keyA == keyB {
		return fmt.Errorf("seed passwords must be 12–72 bytes and distinct tenant keys at least 24 bytes")
	}
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for i, t := range []string{"demo-a", "demo-b"} {
		key := []string{keyA, keyB}[i]
		if _, err = tx.Exec(ctx, `INSERT INTO tenants(id,api_key_hash) VALUES($1,$2) ON CONFLICT(id) DO NOTHING`, t, tokenHash(key)); err != nil {
			return err
		}
	}
	for _, u := range []struct{ ID, Tenant, Email, Role, Pass string }{
		{"admin-a", "demo-a", "admin.a@demo.local", "admin", adminPass},
		{"admin-b", "demo-b", "admin.b@demo.local", "admin", adminPass},
		{"viewer-a", "demo-a", "viewer.a@demo.local", "viewer", viewerPass},
		{"viewer-b", "demo-b", "viewer.b@demo.local", "viewer", viewerPass},
	} {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Pass), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		if _, err = tx.Exec(ctx, `INSERT INTO users(id,tenant,email,role,password_hash) VALUES($1,$2,$3,$4,$5) ON CONFLICT(id) DO NOTHING`, u.ID, u.Tenant, u.Email, u.Role, string(hash)); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
func (p *Postgres) Ping(ctx context.Context) error { return p.Pool.Ping(ctx) }
func (p *Postgres) FindUser(ctx context.Context, email string) (u model.User, err error) {
	err = p.Pool.QueryRow(ctx, `SELECT id,tenant,email,role,password_hash FROM users WHERE email=$1`, email).Scan(&u.ID, &u.Tenant, &u.Email, &u.Role, &u.Hash)
	return
}
func (p *Postgres) CreateSession(ctx context.Context, token, user string, expires time.Time) error {
	_, e := p.Pool.Exec(ctx, `INSERT INTO sessions(token_hash,user_id,expires_at) VALUES($1,$2,$3)`, tokenHash(token), user, expires)
	return e
}
func (p *Postgres) Session(ctx context.Context, token string) (u model.User, err error) {
	err = p.Pool.QueryRow(ctx, `SELECT u.id,u.tenant,u.email,u.role FROM sessions s JOIN users u ON u.id=s.user_id WHERE s.token_hash=$1 AND s.expires_at>now()`, tokenHash(token)).Scan(&u.ID, &u.Tenant, &u.Email, &u.Role)
	return
}
func (p *Postgres) DeleteSession(ctx context.Context, token string) error {
	_, e := p.Pool.Exec(ctx, `DELETE FROM sessions WHERE token_hash=$1`, tokenHash(token))
	return e
}
func (p *Postgres) APIKeyTenant(ctx context.Context, key string) (t string, err error) {
	err = p.Pool.QueryRow(ctx, `SELECT id FROM tenants WHERE api_key_hash=$1`, tokenHash(key)).Scan(&t)
	return
}
func (p *Postgres) Insert(ctx context.Context, events []model.Event) error {
	if len(events) == 0 {
		return nil
	}
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	tenant := events[0].Tenant
	// Serializes per-tenant ingestion and rule evaluation across API/collector processes.
	var rule model.Rule
	if err = tx.QueryRow(ctx, `SELECT alert_enabled,alert_threshold,alert_window FROM tenants WHERE id=$1 FOR UPDATE`, tenant).Scan(&rule.Enabled, &rule.Threshold, &rule.WindowMinutes); err != nil {
		return err
	}
	ips := map[string]bool{}
	for _, e := range events {
		if e.Tenant != tenant {
			return fmt.Errorf("mixed-tenant batch")
		}
		fields, err := json.Marshal(e.Fields)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `INSERT INTO logs(timestamp,tenant,source,event_type,severity,src_ip,username,host,action,fields,raw) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`, e.Timestamp, e.Tenant, e.Source, e.EventType, e.Severity, e.SrcIP, e.User, e.Host, e.Action, fields, []byte(e.Raw))
		if err != nil {
			return err
		}
		if e.EventType == "login_failed" && e.SrcIP != "" {
			ips[e.SrcIP] = true
		}
	}
	if rule.Enabled {
		for ip := range ips {
			_, err = tx.Exec(ctx, `INSERT INTO alerts(tenant,src_ip,event_count)
   SELECT $1,$2,count(*) FROM logs WHERE tenant=$1 AND src_ip=$2 AND event_type='login_failed'
   AND timestamp BETWEEN now()-($3 * interval '1 minute') AND now()
   HAVING count(*) >= $4 AND NOT EXISTS(SELECT 1 FROM alerts WHERE tenant=$1 AND src_ip=$2 AND status='open' AND created_at>now()-($3 * interval '1 minute'))`, tenant, ip, rule.WindowMinutes, rule.Threshold)
			if err != nil {
				return err
			}
		}
	}
	return tx.Commit(ctx)
}
func conditions(f model.Filter) (string, []any) {
	return `tenant=$1 AND timestamp >= $2 AND timestamp <= $3 AND ($4='' OR source=$4) AND ($5='' OR strpos(lower(event_type||' '||src_ip||' '||username||' '||host||' '||fields::text||' '||raw::text),lower($5))>0)`, []any{f.Tenant, f.From, f.To, f.Source, f.Query}
}
func (p *Postgres) Logs(ctx context.Context, f model.Filter) ([]model.Event, int64, error) {
	where, args := conditions(f)
	var total int64
	if err := p.Pool.QueryRow(ctx, `SELECT count(*) FROM logs WHERE `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	args = append(args, f.Limit, f.Offset)
	rows, err := p.Pool.Query(ctx, `SELECT id,timestamp,tenant,source,event_type,severity,src_ip,username,host,action,fields,raw FROM logs WHERE `+where+` ORDER BY timestamp DESC,id DESC LIMIT $6 OFFSET $7`, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	result := []model.Event{}
	for rows.Next() {
		var e model.Event
		var fields, raw []byte
		if err = rows.Scan(&e.ID, &e.Timestamp, &e.Tenant, &e.Source, &e.EventType, &e.Severity, &e.SrcIP, &e.User, &e.Host, &e.Action, &fields, &raw); err != nil {
			return nil, 0, err
		}
		if err = json.Unmarshal(fields, &e.Fields); err != nil {
			return nil, 0, err
		}
		e.Raw = raw
		result = append(result, e)
	}
	return result, total, rows.Err()
}
func (p *Postgres) Stats(ctx context.Context, f model.Filter) (s model.Stats, err error) {
	where, args := conditions(f)
	err = p.Pool.QueryRow(ctx, `SELECT count(*),count(*) FILTER(WHERE severity>=8),count(DISTINCT source) FROM logs WHERE `+where, args...).Scan(&s.Total, &s.Critical, &s.Sources)
	if err != nil {
		return
	}
	err = p.Pool.QueryRow(ctx, `SELECT count(*) FROM alerts WHERE tenant=$1 AND status='open'`, f.Tenant).Scan(&s.Alerts)
	if err != nil {
		return
	}
	group := func(expression string, limit string) ([]model.Count, error) {
		rows, e := p.Pool.Query(ctx, `SELECT `+expression+`,count(*) FROM logs WHERE `+where+` GROUP BY 1 `+limit, args...)
		if e != nil {
			return nil, e
		}
		defer rows.Close()
		out := []model.Count{}
		for rows.Next() {
			var c model.Count
			if e = rows.Scan(&c.Name, &c.Count); e != nil {
				return nil, e
			}
			out = append(out, c)
		}
		return out, rows.Err()
	}
	if s.Timeline, err = group(`to_char(date_trunc('hour',timestamp AT TIME ZONE 'UTC'),'YYYY-MM-DD"T"HH24:00:00"Z"')`, `ORDER BY 1`); err != nil {
		return
	}
	if s.TopIPs, err = group(`coalesce(nullif(src_ip,''),'unknown')`, `ORDER BY 2 DESC,1 LIMIT 5`); err != nil {
		return
	}
	if s.TopUsers, err = group(`coalesce(nullif(username,''),'unknown')`, `ORDER BY 2 DESC,1 LIMIT 5`); err != nil {
		return
	}
	s.TopEvents, err = group(`event_type`, `ORDER BY 2 DESC,1 LIMIT 5`)
	return
}
func (p *Postgres) Alerts(ctx context.Context, tenant string) ([]model.Alert, error) {
	rows, err := p.Pool.Query(ctx, `SELECT id,tenant,src_ip,event_count,created_at,status FROM alerts WHERE tenant=$1 ORDER BY created_at DESC LIMIT 100`, tenant)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []model.Alert{}
	for rows.Next() {
		var a model.Alert
		if err = rows.Scan(&a.ID, &a.Tenant, &a.IP, &a.Count, &a.CreatedAt, &a.Status); err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, rows.Err()
}
func (p *Postgres) Acknowledge(ctx context.Context, tenant string, id int64) (bool, error) {
	r, e := p.Pool.Exec(ctx, `UPDATE alerts SET status='acknowledged' WHERE tenant=$1 AND id=$2`, tenant, id)
	return r.RowsAffected() > 0, e
}
func (p *Postgres) GetRule(ctx context.Context, tenant string) (r model.Rule, e error) {
	e = p.Pool.QueryRow(ctx, `SELECT alert_enabled,alert_threshold,alert_window FROM tenants WHERE id=$1`, tenant).Scan(&r.Enabled, &r.Threshold, &r.WindowMinutes)
	return
}
func (p *Postgres) SetRule(ctx context.Context, tenant string, r model.Rule) error {
	_, e := p.Pool.Exec(ctx, `UPDATE tenants SET alert_enabled=$2,alert_threshold=$3,alert_window=$4 WHERE id=$1`, tenant, r.Enabled, r.Threshold, r.WindowMinutes)
	return e
}

// Retention is based on ingestion time so imported historical samples live at least 7 days.
func (p *Postgres) Cleanup(ctx context.Context, days int) error {
	if days < 7 {
		return fmt.Errorf("retention must be >=7 days")
	}
	_, e := p.Pool.Exec(ctx, `DELETE FROM logs WHERE ingested_at < now()-($1 * interval '1 day')`, days)
	if e != nil {
		return e
	}
	_, e = p.Pool.Exec(ctx, `DELETE FROM sessions WHERE expires_at<now()`)
	return e
}
