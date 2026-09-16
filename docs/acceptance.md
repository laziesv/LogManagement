# Assignment checklist

Implementation is separate from acceptance verification. See verification.md for executed checks.

## Implemented locally

- [x] React + Go Fiber + PostgreSQL source code
- [x] HTTP JSON and file batch ingestion, Go UDP/TCP Syslog receiver
- [x] Samples for API, firewall, network, CrowdStrike, AWS, M365, AD
- [x] Shared normalization, raw preservation, validation, transactional batch
- [x] Search, source/time filters, pagination and event details
- [x] Dashboard totals, timeline, Top IP/User/EventType
- [x] Failed-login rule, alert history and acknowledgement
- [x] Admin/Viewer sessions and tenant-scoped reads/writes
- [x] Tenant enforcement from session/API key instead of user-selected UI filter
- [x] Minimum seven-day retention worker
- [x] Compose packaging, random .env generator, scripts and documentation
- [x] Unit/HTTP authorization tests and running-stack smoke script

## Must verify before submission

- [ ] Clean Docker Compose build/start on Ubuntu-sized environment
- [ ] Syslog appears in UI within one minute
- [ ] HTTP, uploaded JSON and simulator produce searchable normalized data
- [ ] Import AWS, M365 and AD samples and inspect fields/raw
- [ ] Verify two tenants with distinct data and Viewer restrictions against real DB
- [ ] Confirm alert creation/cooldown/acknowledgement against real DB
- [ ] Confirm retention with old ingested_at records in an isolated test database
- [ ] Deploy cloud HTTPS and provide evaluator URL/account
- [ ] Record the 30-minute demonstration and schedule interview within deadline

## Optional later

- [ ] CI/CD, load test results, metrics/tracing
- [ ] Dedicated tables/indexes per tenant or PostgreSQL RLS defense in depth
- [ ] Queue/retry/dedup, enrichment and additional provider field mappings
- [ ] Postman/Insomnia collection
