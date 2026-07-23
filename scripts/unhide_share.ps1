param(
    [Parameter(Mandatory = $true)]
    [string]$Path
)

if (-not (Test-Path -LiteralPath $Path)) {
    throw "Path not found: $Path"
}

$full = (Resolve-Path -LiteralPath $Path).Path
Write-Host "Processing: $full"

Get-ChildItem -LiteralPath $full -Force -Recurse -ErrorAction Stop | ForEach-Object {
    try {
        $item = $_
        if (($item.Attributes -band [IO.FileAttributes]::Hidden) -or ($item.Attributes -band [IO.FileAttributes]::System)) {
            $newAttrs = $item.Attributes -band (-bnot [IO.FileAttributes]::Hidden) -band (-bnot [IO.FileAttributes]::System)
            $item.Attributes = $newAttrs
            Write-Host "Updated: $($item.FullName)"
        }
    } catch {
        Write-Warning "Failed: $($_.Exception.Message)"
    }
}

try {
    $root = Get-Item -LiteralPath $full -Force
    if (($root.Attributes -band [IO.FileAttributes]::Hidden) -or ($root.Attributes -band [IO.FileAttributes]::System)) {
        $root.Attributes = $root.Attributes -band (-bnot [IO.FileAttributes]::Hidden) -band (-bnot [IO.FileAttributes]::System)
        Write-Host "Updated root: $full"
    }
} catch {
    Write-Warning "Failed root update: $($_.Exception.Message)"
}

Write-Host "Done."
