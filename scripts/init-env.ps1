$ErrorActionPreference = 'Stop'
$projectDir = Split-Path -Parent $PSScriptRoot
$envPath = Join-Path $projectDir '.env'
if (Test-Path -LiteralPath $envPath) { Write-Host '.env already exists; preserved.'; exit 0 }
function New-Secret([int]$Length = 24) {
  $bytes = New-Object byte[] $Length
  $rng = [System.Security.Cryptography.RandomNumberGenerator]::Create()
  try { $rng.GetBytes($bytes) } finally { $rng.Dispose() }
  return ([BitConverter]::ToString($bytes)).Replace('-', '').ToLowerInvariant()
}
$lines = @(
  "DB_PASSWORD=$(New-Secret 24)",
  "ADMIN_PASSWORD=$(New-Secret 12)",
  "VIEWER_PASSWORD=$(New-Secret 12)",
  "API_KEY_A=$(New-Secret 32)",
  "API_KEY_B=$(New-Secret 32)",
  'APP_ORIGIN=http://localhost:8080',
  'COOKIE_SECURE=false',
  'RETENTION_DAYS=7',
  'HTTP_BIND=127.0.0.1',
  'SYSLOG_BIND=127.0.0.1'
)
[IO.File]::WriteAllLines($envPath, $lines)
Write-Host 'Created .env with random credentials. Read ADMIN_PASSWORD in .env to sign in as admin.a@demo.local.'
