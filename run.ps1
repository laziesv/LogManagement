$ErrorActionPreference = 'Stop'
Set-Location -LiteralPath $PSScriptRoot
& "$PSScriptRoot/scripts/init-env.ps1"
docker compose up --build -d
if ($LASTEXITCODE -ne 0) { throw 'Docker Compose failed. Make sure Docker Engine is running.' }
