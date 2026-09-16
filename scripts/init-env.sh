#!/usr/bin/env sh
set -eu
cd "$(dirname "$0")/.."
if [ -f .env ]; then echo '.env already exists; preserved.'; exit 0; fi
umask 077
secret() { openssl rand -hex "$1"; }
{
 printf 'DB_HOST=postgres\n'
 printf 'DB_PORT=5432\n'
 printf 'DB_NAME=logmanagement\n'
 printf 'DB_USER=logmanagement\n'
 printf 'DB_PASSWORD=%s\n' "$(secret 24)"
 printf 'ADMIN_PASSWORD=%s\n' "$(secret 12)"
 printf 'VIEWER_PASSWORD=%s\n' "$(secret 12)"
 printf 'API_KEY_A=%s\n' "$(secret 32)"
 printf 'API_KEY_B=%s\n' "$(secret 32)"
 printf 'APP_ORIGIN=http://localhost:8080\nCOOKIE_SECURE=false\nRETENTION_DAYS=7\nHTTP_BIND=127.0.0.1\nSYSLOG_BIND=127.0.0.1\nSYSLOG_TENANT=demo-a\n'
} > .env
echo 'Created .env. Read ADMIN_PASSWORD to sign in as admin.a@demo.local.'
