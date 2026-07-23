param(
    [string]$Path = "\\RED01\QUIROFANO",
    [int]$IntervalMinutes = 5,
    [switch]$Once
)

# Ejemplos de uso:
# Ruta UNC desde otro equipo:
# .\unhide_share.ps1 -Path "\\RED01\QUIROFANO"
#
# Ruta local cuando el script corre en el propio servidor RED01:
# .\unhide_share.ps1 -Path "D:\RED01\QUIROFANO"
#
# Si lo van a dejar como tarea programada en RED01, usa la ruta física real
# del servidor y descomenta o reemplaza el valor de -Path según corresponda.

$ErrorActionPreference = "Stop"

function Clear-HiddenAttributes {
    param(
        [Parameter(Mandatory = $true)]
        [string]$TargetPath
    )

    if (-not (Test-Path -LiteralPath $TargetPath)) {
        Write-Warning "No existe: $TargetPath"
        return
    }

    $items = Get-ChildItem -LiteralPath $TargetPath -Force -Recurse -ErrorAction SilentlyContinue
    foreach ($item in $items) {
        try {
            $attrs = $item.Attributes
            if (($attrs -band [IO.FileAttributes]::Hidden) -or ($attrs -band [IO.FileAttributes]::System)) {
                $item.Attributes = $attrs -band (-bnot [IO.FileAttributes]::Hidden) -band (-bnot [IO.FileAttributes]::System)
                Write-Host "Updated: $($item.FullName)"
            }
        } catch {
            Write-Warning "Failed: $($item.FullName) :: $($_.Exception.Message)"
        }
    }

    try {
        $root = Get-Item -LiteralPath $TargetPath -Force
        if (($root.Attributes -band [IO.FileAttributes]::Hidden) -or ($root.Attributes -band [IO.FileAttributes]::System)) {
            $root.Attributes = $root.Attributes -band (-bnot [IO.FileAttributes]::Hidden) -band (-bnot [IO.FileAttributes]::System)
            Write-Host "Updated root: $TargetPath"
        }
    } catch {
        Write-Warning "Failed root: $TargetPath :: $($_.Exception.Message)"
    }
}

function Start-UnhideLoop {
    param(
        [Parameter(Mandatory = $true)]
        [string]$WatchPath,
        [Parameter(Mandatory = $true)]
        [int]$EveryMinutes
    )

    $interval = [TimeSpan]::FromMinutes($EveryMinutes)
    if ($interval.TotalSeconds -lt 30) {
        $interval = [TimeSpan]::FromMinutes(1)
    }

    Write-Host "Watching $WatchPath every $($interval.TotalMinutes) minute(s). Press Ctrl+C to stop."
    while ($true) {
        Clear-HiddenAttributes -TargetPath $WatchPath
        Start-Sleep -Seconds [int]$interval.TotalSeconds
    }
}

$fullPath = (Resolve-Path -LiteralPath $Path).Path

Clear-HiddenAttributes -TargetPath $fullPath

if ($Once) {
    Write-Host "Done."
    exit 0
}

Start-UnhideLoop -WatchPath $fullPath -EveryMinutes $IntervalMinutes
