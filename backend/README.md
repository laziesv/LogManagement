# Go Fiber backend

Go 1.25+. Run `go test ./...` and `go vet ./...`.

`cmd/server` only loads configuration, receives OS signals and calls `bootstrap.Run`.

- `internal/config`: environment defaults and validation.
- `internal/bootstrap`: database initialization, dependency wiring, HTTP startup and graceful shutdown.
- `internal/collector`: UDP/TCP Syslog listeners and connection lifecycle.
- `internal/retention`: immediate/hourly cleanup worker.

- `internal/router`: Fiber server setup and the endpoint table in `routes.go`.
- `internal/handler`: named HTTP handlers, grouped by auth, ingest, logs, alerts and health.
- `internal/middleware`: authentication, Admin checks, ingest authorization, Origin checks and error handling.
- `internal/model`: shared data structures and JSON contracts only; no Fiber or SQL dependencies.
- `internal/normalize`: provider mapping, Syslog normalization and JSON batch decoding.
- `internal/repository`: Store interface, PostgreSQL implementation and embedded schema.

Dependency flow: `bootstrap → router → handler/middleware → repository → model`; handlers and collector use `normalize → model`. Routes keep the existing URLs and middleware ordering. Tests under `router` exercise public HTTP behavior; normalization tests live with `normalize`.

Use root Docker Compose for startup, or set the environment documented in `../docs/setup_appliance.md` and run `go run ./cmd/server`.

The schema is idempotent initial bootstrap, not a versioned migration framework. Future schema changes need explicit migrations. Startup inserts missing demo accounts; it never silently rotates existing credentials.
