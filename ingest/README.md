# Ingestion

The collector is implemented in Go inside `backend/internal/collector/syslog.go`, wired by `backend/internal/bootstrap/run.go` and calls the same normalizer as Fiber (`backend/internal/normalize/normalize.go`). It runs in the backend container for simple appliance packaging.

Channels: HTTP JSON POST, multipart JSON file, Syslog UDP/TCP, Python simulator. Sources: firewall/network/api/crowdstrike/aws/m365/ad. For Syslog the tenant comes from SYSLOG_TENANT, not device-supplied text. JSON ingestion uses the authenticated API key or admin session.

See `samples/` for senders and payloads, and `docs/architecture.md` for supported formats and limitations.
