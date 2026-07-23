//go:build windows

package app

import (
	"fmt"
	"os"
	"syscall"
)

type windowsAccessReport struct {
	Path       string
	Exists     bool
	Accessible bool
	Hidden     bool
	Reason     string
}

func inspectWindowsPath(path string) windowsAccessReport {
	report := windowsAccessReport{Path: path}
	st, err := os.Stat(path)
	if err != nil {
		report.Reason = classifyFileError(err)
		return report
	}
	report.Exists = true
	report.Hidden = windowsHidden(path, st)
	f, err := os.Open(path)
	if err != nil {
		report.Reason = classifyFileError(err)
		if report.Hidden {
			report.Reason = appendReason(report.Reason, "oculto")
		}
		return report
	}
	_ = f.Close()
	report.Accessible = true
	if report.Hidden {
		report.Reason = "oculto"
	}
	return report
}

func unhideWindowsPath(path string) (string, error) {
	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return "", err
	}
	attrs, err := syscall.GetFileAttributes(p)
	if err != nil {
		return "", err
	}
	attrs &^= syscall.FILE_ATTRIBUTE_HIDDEN
	attrs &^= syscall.FILE_ATTRIBUTE_SYSTEM
	if err := syscall.SetFileAttributes(p, attrs); err != nil {
		return "", err
	}
	return "atributos Hidden/System removidos", nil
}

func windowsHidden(path string, st os.FileInfo) bool {
	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false
	}
	attrs, err := syscall.GetFileAttributes(p)
	if err != nil {
		return false
	}
	return attrs&syscall.FILE_ATTRIBUTE_HIDDEN != 0
}

func windowsAccessSupported() bool {
	return true
}

func windowsAccessHelp() string {
	return "Diagnosticar acceso"
}

func windowsFixHelp() string {
	return "Quitar oculto"
}

func windowsActionHelp() string {
	return "Diagnóstico y corrección nativa"
}

func windowsPathActionError() error {
	return fmt.Errorf("no se pudo operar sobre la ruta")
}

func windowsPathActionLabel() string {
	return "Windows"
}
