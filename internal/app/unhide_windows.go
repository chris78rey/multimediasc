//go:build windows

package app

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
)

// ClearHiddenAttributes removes Hidden and System attributes from the target path
// and everything beneath it.
func ClearHiddenAttributes(targetPath string) error {
	info, err := os.Stat(targetPath)
	if err != nil {
		return err
	}

	if err := clearHiddenAttributes(targetPath); err != nil {
		return err
	}

	if info.IsDir() {
		return filepath.WalkDir(targetPath, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if path == targetPath {
				return nil
			}
			return clearHiddenAttributes(path)
		})
	}

	return nil
}

func clearHiddenAttributes(path string) error {
	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	attrs, err := syscall.GetFileAttributes(p)
	if err != nil {
		return err
	}

	cleaned := attrs &^ (syscall.FILE_ATTRIBUTE_HIDDEN | syscall.FILE_ATTRIBUTE_SYSTEM)
	if cleaned == attrs {
		return nil
	}

	if err := syscall.SetFileAttributes(p, cleaned); err != nil {
		return fmt.Errorf("set attributes for %s: %w", path, err)
	}

	return nil
}
