# Appliance setup

Target: Ubuntu 22.04+, 4 vCPU, 8 GB RAM, 40 GB disk, Docker Engine + Compose plugin. Windows development: Docker Desktop in Linux container mode.

## Start

Windows PowerShell from the project root:

```powershell
.\run.ps1
```

Ubuntu/Linux:

```sh
sh run.sh
```

The script generates a random `.env` only if it does not exist and builds/starts PostgreSQL, backend and frontend. First start needs Internet for container images and dependencies.

Visit http://localhost:8080. Read ADMIN_PASSWORD in .env for admin@demo.local; Viewer accounts use VIEWER_PASSWORD. Do not include .env in the submitted repository.

```sh
docker compose ps
docker compose logs --tail=100 backend
curl http://localhost:8080/api/health
```

Services start in health-check order. Data lives in the named postgres_data volume and survives `docker compose down` (without `-v`). Never delete the volume to fix an issue unless its data is disposable.

## Ingest examples

UI: Data sources → Send demo batch, or Import logs → samples/events.json.

Syslog:

```sh
python samples/send_syslog.py
python samples/send_syslog.py --tcp
```

HTTP simulator (set LOG_API_KEY to API_KEY_A/B from .env first):

```sh
python samples/post_logs.py
python samples/post_logs.py --alert
python samples/post_logs.py --file samples/aws.json
```

File examples deliberately omit tenant so credentials determine ownership. Omitted timestamps use ingestion time. Files with historical timestamps require a matching UI Custom range (up to 31 days).

## Network exposure

Defaults bind HTTP 8080 and Syslog 5514 UDP/TCP to 127.0.0.1. PostgreSQL and the backend HTTP port are not published. For trusted LAN testing, set HTTP_BIND and/or SYSLOG_BIND to the specific private interface and APP_ORIGIN to the actual browser origin. Firewall Syslog to known devices; never expose unauthenticated Syslog broadly to the Internet.

The host Syslog port is 5514 to avoid privileged-port conflicts. If port 514 is required, change only the host port in Compose and use `--port 514` in the sender.

## Local development without rebuilding containers

Use a local PostgreSQL instance or start the Compose database with a development-only loopback port mapping. Set DATABASE_URL, ADMIN_PASSWORD, VIEWER_PASSWORD, API_KEY_A, API_KEY_B, APP_ORIGIN and COOKIE_SECURE in the backend process environment. Use APP_ORIGIN=http://127.0.0.1:5173 for the default Vite command and COOKIE_SECURE=false.

```sh
cd backend
go run ./cmd/server
# In another terminal:
cd frontend
npm ci
npm run dev
```

## Troubleshooting

- Docker pipe/daemon unavailable: open Docker Desktop and wait for Engine running. If Docker requests WSL setup, license acceptance or a reboot, complete it on the host.
- Login invalid: read the initial .env password. Seeds do not update existing users or keys on restart.
- 403 on writes: match APP_ORIGIN to the browser URL exactly, including scheme and port. Cloud requires COOKIE_SECURE=true.
- No logs: check source/time/search filters; collector tenant is demo-a. Refresh is manual.
- API unavailable: inspect `docker compose logs backend postgres`; UI does not silently switch to mock data.
