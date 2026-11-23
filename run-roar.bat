@echo off
setlocal
cd /d C:\Users\jeeva\ROAR-clean

echo Building ROAR...
where go >nul 2>&1
if errorlevel 1 (
  echo Go was not found in PATH. Use the full Go installation or run from a Go terminal.
  pause
  exit /b 1
)
go build -o roar.exe ./cmd/roar
if errorlevel 1 (
  echo Build failed.
  pause
  exit /b 1
)

echo Starting ROAR at http://127.0.0.1:8000/ui/
echo Keep this window open. Press Ctrl+C to stop.
roar.exe api

echo ROAR stopped. Check the message above.
pause
