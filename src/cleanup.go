package main

import (
	"os"
	"path/filepath"
	"strings"
)

// clearTemporaryFiles removes entries from the current user's temporary
// directory and from %WINDIR%\Temp. Failures are intentionally ignored:
// locked files and system-protected entries are common on Windows.
func clearTemporaryFiles() {
	seen := map[string]bool{}
	if p := os.TempDir(); p != "" {
		clearDirectoryContents(p, seen)
	}
	if winDir := os.Getenv("WINDIR"); winDir != "" {
		clearDirectoryContents(filepath.Join(winDir, "Temp"), seen)
	}
}

func clearDirectoryContents(dir string, seen map[string]bool) {
	clean := filepath.Clean(dir)
	key := strings.ToLower(clean)
	if clean == "." || clean == string(filepath.Separator) || seen[key] {
		return
	}
	seen[key] = true

	entries, err := os.ReadDir(clean)
	if err != nil {
		return
	}
	exeDir, _ := filepath.Abs(executableDir())
	exeDirKey := strings.ToLower(filepath.Clean(exeDir))
	for _, entry := range entries {
		target := filepath.Join(clean, entry.Name())
		targetAbs, _ := filepath.Abs(target)
		targetKey := strings.ToLower(filepath.Clean(targetAbs))
		// Never remove the directory that contains the running MiniBin instance.
		if exeDirKey == targetKey || strings.HasPrefix(exeDirKey, targetKey+string(filepath.Separator)) {
			continue
		}
		_ = os.RemoveAll(target)
	}
}
