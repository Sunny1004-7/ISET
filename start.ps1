# Start frontend and backend services
$rootPath = if ($PSScriptRoot) { $PSScriptRoot } else { Get-Location }

Write-Host "Starting services..." -ForegroundColor Cyan

# Start backend
Write-Host "Starting backend..." -ForegroundColor Green
$backendPath = Join-Path $rootPath "backend"
if (Test-Path $backendPath) {
    $cmd = "cd '$backendPath'; ./backend"
    Start-Process powershell -ArgumentList "-NoExit", "-Command", $cmd
    Write-Host "Backend window opened" -ForegroundColor Green
}

Start-Sleep -Milliseconds 500

# Start frontend
Write-Host "Starting frontend..." -ForegroundColor Green
$frontendPath = Join-Path $rootPath "ISET-frontend"
if (Test-Path $frontendPath) {
    $cmd = "cd '$frontendPath'; pnpm install; pnpm serve"
    Start-Process powershell -ArgumentList "-NoExit", "-Command", $cmd
    Write-Host "Frontend window opened" -ForegroundColor Green
}

Write-Host "Backend: http://localhost:8010" -ForegroundColor Yellow
Write-Host "Frontend: http://localhost:8021" -ForegroundColor Yellow
