# Optional Nginx access logs

The frontend Nginx service keeps access logs local by default, so refreshing the web UI does not create new application logs. The Nginx access format still omits query strings, headers, cookies and request bodies.

To demonstrate infrastructure access-log ingestion, temporarily add this line inside the `server` block in `frontend/nginx.conf`, next to the local `access_log` line:

```nginx
access_log syslog:server=backend:5514,facility=local7,tag=logdesk_access,severity=info logdesk_access;
```

Then rebuild the frontend container and reload http://localhost:8080 to generate actual access traffic. Sign in as Admin or Viewer A, open Log explorer, select source `network`, search `nginx_access`, and press Refresh.

Inspect an event: `fields` includes `origin: live`, `vendor: nginx`, HTTP method, URL path, status code, response time in seconds, bytes sent, and the peer IP observed by Nginx. `raw` preserves the received Syslog message.

Normal responses have severity 3, HTTP 4xx severity 5, and HTTP 5xx severity 8. These access events do not trigger the failed-login alert rule: an HTTP error is not necessarily a failed login.

The source is `network` to match the existing source schema. The event type `nginx_access` and `origin: live` distinguish it from synthetic samples. Infrastructure access logs belong to the collector tenant, regardless of which application user visits the website. Viewer B cannot read demo-a logs. Docker may present a gateway IP as the client address.

API routes and `/ingest` should stay excluded so reading or ingesting logs does not create recursive collection traffic. URL paths are still stored: avoid putting secrets in path segments. Container stdout uses the same sanitized access format.

UDP collection is best effort and has no retry queue: traffic during collector downtime can be lost. When replacing the backend container independently, restart the frontend so Nginx resolves its current address.

## Verify

With Syslog forwarding enabled, run `python tests/nginx_live.py` against the running local stack. It generates one actual missing-asset request and confirms storage, field normalization, query omission and tenant isolation. The event remains available for inspection.

References: [Nginx Syslog](https://nginx.org/en/docs/syslog.html), [access log formatting](https://nginx.org/en/docs/http/ngx_http_log_module.html).
