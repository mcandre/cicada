//go:build android

package cicada

import (
	"runtime"
)

var EnvironmentIsLinux bool = true

// RecognizeOs identifies the environment,
// as an endoflife.date product name.
func RecognizeOs() (*string, error) {
	goos := runtime.GOOS
	return &goos, nil
}
