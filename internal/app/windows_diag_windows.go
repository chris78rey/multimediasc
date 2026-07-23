//go:build windows

package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type windowsFileDiagnostic struct {
	Path       string `json:"path"`
	Exists     bool   `json:"exists"`
	Accessible bool   `json:"accessible"`
	Hidden     bool   `json:"hidden"`
	Reason     string `json:"reason"`
}

func diagnoseWindowsPath(path string) (windowsFileDiagnostic, error) {
	script := fmt.Sprintf(`
$ErrorActionPreference = 'Stop'
$p = '%s'
$item = Get-Item -LiteralPath $p -Force -ErrorAction SilentlyContinue
if ($null -eq $item) {
  [pscustomobject]@{ path=$p; exists=$false; accessible=$false; hidden=$false; reason='no existe o no es visible con la identidad actual' } | ConvertTo-Json -Compress
  exit 0
}
$hidden = [bool]($item.Attributes -band [System.IO.FileAttributes]::Hidden)
$accessible = $true
$reason = 'acceso permitido'
try {
  $null = Get-Content -LiteralPath $p -TotalCount 1 -ErrorAction Stop
} catch {
  $accessible = $false
  $reason = $_.Exception.Message
}
[pscustomobject]@{ path=$p; exists=$true; accessible=$accessible; hidden=$hidden; reason=$reason } | ConvertTo-Json -Compress
`, escapePowerShellSingleQuotes(path))

	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", script)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return windowsFileDiagnostic{}, fmt.Errorf("%v: %s", err, strings.TrimSpace(out.String()))
	}
	var diag windowsFileDiagnostic
	if err := json.Unmarshal(bytes.TrimSpace(out.Bytes()), &diag); err != nil {
		return windowsFileDiagnostic{}, fmt.Errorf("no se pudo interpretar el diagnóstico de PowerShell: %w", err)
	}
	return diag, nil
}

func fixWindowsHidden(path string) (string, error) {
	script := fmt.Sprintf(`
$ErrorActionPreference = 'Stop'
$p = '%s'
if (-not (Test-Path -LiteralPath $p -Force)) {
  Write-Output 'no existe'
  exit 0
}
$item = Get-Item -LiteralPath $p -Force
$item.Attributes = $item.Attributes -band (-bnot [System.IO.FileAttributes]::Hidden) -band (-bnot [System.IO.FileAttributes]::System)
Write-Output 'oculto removido'
`, escapePowerShellSingleQuotes(path))

	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", script)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v: %s", err, strings.TrimSpace(out.String()))
	}
	return strings.TrimSpace(out.String()), nil
}

func escapePowerShellSingleQuotes(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
