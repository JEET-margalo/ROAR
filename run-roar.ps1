$ErrorActionPreference = 'Stop'
Set-Location $PSScriptRoot

Write-Host 'Building ROAR...' -ForegroundColor Cyan
& go build -o roar.exe ./cmd/roar
if ($LASTEXITCODE -ne 0) { throw 'Build failed' }

$existing = Get-Process roar -ErrorAction SilentlyContinue
if ($existing) {
  Write-Host 'Stopping existing ROAR process...' -ForegroundColor Yellow
  $existing | Stop-Process -Force
  Start-Sleep -Milliseconds 500
}

Write-Host 'ROAR is starting at http://127.0.0.1:8000/ui/' -ForegroundColor Green
Write-Host 'Keep this window open while using the website.' -ForegroundColor Yellow
& .\roar.exe api
