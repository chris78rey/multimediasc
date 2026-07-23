//go:build !windows

package app

import (
	"os"
	"strings"
)

func isHidden(path string, st os.FileInfo) bool {
	base := strings.TrimSpace(st.Name())
	return strings.HasPrefix(base, ".")
}
