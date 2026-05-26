#!/usr/bin/env pwsh
# ExchangeOS — PowerShell mirror para Windows (FX-XOS-001).
# Delega para Task quando disponivel, fallback para Go nativo.

param(
    [Parameter(Mandatory=$true, Position=0)]
    [ValidateSet('install','build','test','lint','fmt','sec','db-up','db-migrate','compose-up','compose-down','dash','clean','help')]
    [string]$Command
)

$ErrorActionPreference = 'Stop'
$Root = Split-Path -Parent $PSScriptRoot

function Invoke-Task {
    param([string]$Target)
    if (Get-Command task -ErrorAction SilentlyContinue) {
        & task $Target
    } else {
        Write-Error "Task runner not installed. Install via: scoop install task (or chocolatey/winget)"
        exit 1
    }
}

switch ($Command) {
    'install'      { Invoke-Task 'install' }
    'build'        { Invoke-Task 'build' }
    'test'         { Invoke-Task 'test' }
    'lint'         { Invoke-Task 'lint' }
    'fmt'          { Invoke-Task 'fmt' }
    'sec'          { Invoke-Task 'sec:secrets'; Invoke-Task 'sec:trivy'; Invoke-Task 'sec:govulncheck' }
    'db-up'        { Invoke-Task 'db:up' }
    'db-migrate'   { Invoke-Task 'db:migrate' }
    'compose-up'   { Invoke-Task 'compose:up' }
    'compose-down' { Invoke-Task 'compose:down' }
    'dash'         { Invoke-Task 'dash-update' }
    'clean'        { Invoke-Task 'clean' }
    'help'         { Invoke-Task 'default' }
}
