param(
    [Parameter(Position = 0)]
    [string]$Command = "help",

    [Parameter(Position = 1)]
    [string]$Arg
)

$ErrorActionPreference = "Stop"

$root = $PSScriptRoot
$dbInitDir = Join-Path $root "backend/db/init"

function Show-Help {
    Write-Host "Yardly helper commands:"
    Write-Host "  .\commands.ps1 help"
    Write-Host "  .\commands.ps1 db-up"
    Write-Host "  .\commands.ps1 db-down"
    Write-Host "  .\commands.ps1 db-shell"
    Write-Host "  .\commands.ps1 backend"
    Write-Host "  .\commands.ps1 frontend"
    Write-Host "  .\commands.ps1 dev"
    Write-Host "  .\commands.ps1 migration:new <name>"
    Write-Host "  .\commands.ps1 migration:apply <file.sql>"
}

function Ensure-Docker {
    if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
        throw "docker is not installed or not in PATH"
    }
}

function Start-Db {
    Ensure-Docker
    Push-Location $root
    try {
        docker compose up -d db
    }
    finally {
        Pop-Location
    }
}

function Stop-Db {
    Ensure-Docker
    Push-Location $root
    try {
        docker compose stop db
    }
    finally {
        Pop-Location
    }
}

function Open-DbShell {
    Ensure-Docker
    docker exec -it yardly-db psql -U yardly -d yardly
}

function Run-Backend {
    Push-Location (Join-Path $root "backend")
    try {
        go run ./cmd/api
    }
    finally {
        Pop-Location
    }
}

function Run-Frontend {
    Push-Location (Join-Path $root "frontend")
    try {
        npm run dev
    }
    finally {
        Pop-Location
    }
}

function Run-Dev {
    Start-Db

    $backendPath = Join-Path $root "backend"
    $frontendPath = Join-Path $root "frontend"

    Start-Process powershell -ArgumentList "-NoExit", "-Command", "Set-Location '$backendPath'; go run ./cmd/api"
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "Set-Location '$frontendPath'; npm run dev"

    Write-Host "Opened 2 terminals: backend + frontend"
}

function New-Migration {
    param([Parameter(Mandatory = $true)][string]$Name)

    if (-not (Test-Path $dbInitDir)) {
        throw "Migration dir not found: $dbInitDir"
    }

    $clean = ($Name.Trim() -replace "\s+", "_" -replace "[^a-zA-Z0-9_]", "").ToLower()
    if ([string]::IsNullOrWhiteSpace($clean)) {
        throw "Migration name is empty after cleanup"
    }

    $existing = Get-ChildItem -Path $dbInitDir -Filter "*.sql" | Select-Object -ExpandProperty BaseName
    $max = 0

    foreach ($item in $existing) {
        if ($item -match "^(\d{3})_") {
            $n = [int]$Matches[1]
            if ($n -gt $max) {
                $max = $n
            }
        }
    }

    $next = $max + 1
    $fileName = "{0:D3}_{1}.sql" -f $next, $clean
    $fullPath = Join-Path $dbInitDir $fileName

    @(
        "-- +migrate Up",
        "",
        "",
        "-- +migrate Down",
        ""
    ) | Set-Content -Path $fullPath -Encoding UTF8

    Write-Host "Created migration: backend/db/init/$fileName"
}

function Apply-Migration {
    param([Parameter(Mandatory = $true)][string]$File)

    Ensure-Docker

    if (-not ($File.ToLower().EndsWith(".sql"))) {
        throw "Pass sql file name, example: 023_add_index.sql"
    }

    $fullPath = Join-Path $dbInitDir $File
    if (-not (Test-Path $fullPath)) {
        throw "File not found: backend/db/init/$File"
    }

    docker exec -it yardly-db psql -v ON_ERROR_STOP=1 -U yardly -d yardly -f "/docker-entrypoint-initdb.d/$File"
}

switch ($Command.ToLower()) {
    "help" { Show-Help }
    "db-up" { Start-Db }
    "db-down" { Stop-Db }
    "db-shell" { Open-DbShell }
    "backend" { Run-Backend }
    "frontend" { Run-Frontend }
    "dev" { Run-Dev }
    "migration:new" {
        if (-not $Arg) { throw "Usage: .\commands.ps1 migration:new <name>" }
        New-Migration -Name $Arg
    }
    "migration:apply" {
        if (-not $Arg) { throw "Usage: .\commands.ps1 migration:apply <file.sql>" }
        Apply-Migration -File $Arg
    }
    default {
        Write-Host "Unknown command: $Command"
        Show-Help
        exit 1
    }
}
