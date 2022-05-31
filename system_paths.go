package cicada

import (
	"path/filepath"
)

// SystemPaths documents common stock executable directories,
// which are eligible for skipping in quiet mode.
var SystemPaths = []string{
	"/bin",
	"/usr/bin",
	"/usr/sbin",
	"/usr/share/bin",
	"/sbin",
	"c:\\Windows",
	"c:\\Windows\\system32",
	"c:\\Windows\\System32\\Wbem",
	"/mnt/c/Windows",
	"/mnt/c/Windows/system32",
	"/mnt/c/Windows/System32/Wbem",
}

// IsSystemExecutable returns true
// when the given executable path is a child of a system directory
// in SystemPaths.
//
// Otherwise, returns false.
func IsSystemExecutable(executablePath string) bool {
	for _, systemPath := range SystemPaths {
		if filepath.Dir(executablePath) == systemPath {
			return true
		}
	}

	return false
}
