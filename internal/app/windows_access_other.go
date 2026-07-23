//go:build !windows

package app

import "fmt"

type windowsAccessReport struct {
	Path       string
	Exists     bool
	Accessible bool
	Hidden     bool
	Reason     string
}

func inspectWindowsPath(path string) windowsAccessReport {
	return windowsAccessReport{Path: path}
}

func unhideWindowsPath(path string) (string, error) {
	return "", fmt.Errorf("solo disponible en Windows")
}

func windowsAccessSupported() bool {
	return false
}

func windowsAccessHelp() string {
	return "No disponible"
}

func windowsFixHelp() string {
	return "No disponible"
}

func windowsActionHelp() string {
	return "Solo Windows"
}

func windowsPathActionError() error {
	return fmt.Errorf("no disponible fuera de Windows")
}

func windowsPathActionLabel() string {
	return "no-op"
}
