# SaaS / Cloud deployment

Status: deployment instructions and HTTPS template are provided; no cloud resource, DNS record or public URL has been created.

Use a VM with Ubuntu 22.04+, 4 vCPU, 8 GB RAM, 40 GB disk, Docker/Compose, host Nginx and a domain pointing to the VM. Keep the database and backend private.

1. Copy the repository to the VM; run `sh scripts/init-env.sh`.
2. Set APP_ORIGIN=https://YOUR_DOMAIN and COOKIE_SECURE=true in .env. Keep HTTP_BIND=127.0.0.1 and SYSLOG_BIND=127.0.0.1.
3. Run `docker compose up --build -d` and verify localhost:8080/api/health.
4. Obtain a TLS certificate for your domain using your chosen ACME client/issuer. Install host Nginx with the actual domain/certificate paths from deploy/nginx-tls.conf.example. Validate with `sudo nginx -t` before reloading.
5. Expose 80/443 through the cloud firewall; restrict SSH to your trusted source. Do not expose database port 5432 or public Syslog.
6. Visit https://YOUR_DOMAIN, log in and send a sample. Verify the session cookie has Secure and HttpOnly attributes. Check Viewer B cannot see demo-a logs.
7. Run `python tests/smoke.py` with the same root .env and HTTPS APP_ORIGIN. Standard CA validation stays enabled.
8. Share only the public URL and intended demo-account credentials directly with the evaluator. Do not put secrets in Git or documentation.

For the assignment's self-signed option, provide certificate trust/import instructions to the evaluator. Do not disable TLS validation in the application or scripts. A publicly trusted certificate is easier for the evaluator.

Use backups for the PostgreSQL volume, restrict server access, and monitor disk capacity. This demo uses shared-table tenant isolation at the application layer and has not been load tested for production scale.
