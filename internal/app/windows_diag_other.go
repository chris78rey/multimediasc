//go:build !windows

package app

import "fmt"

type windowsFileDiagnostic struct {
	Path       string
	Exists     bool
	Accessible bool
	Hidden     bool
	Reason     string
}

func diagnoseWindowsPath(path string) (windowsFileDiagnostic, error) {
	return windowsFileDiagnostic{}, fmt.Errorf("diagnóstico PowerShell solo disponible en Windows")
}

func fixWindowsHidden(path string) (string, error) {
	return "", fmt.Errorf("corrección PowerShell solo disponible en Windows")
}
