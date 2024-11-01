//go:build android

package cicada

import (
	"runtime"
)

// EnvironmentIsLinux checks whether the current platform is Linux.
var EnvironmentIsLinux = true

// RecognizeOs identifies the environment,
// as an endoflife.date product name.
func RecognizeOs() (*string, error) {
	goos := runtime.GOOS
	return &goos, nil
}
