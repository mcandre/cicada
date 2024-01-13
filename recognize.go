//go:build !darwin && !linux

package cicada

import (
	"runtime"
)

var EnvironmentIsLinux bool = false

// RecognizeOs identifies the environment,
// as an endoflife.date product name.
func RecognizeOs() (*string, error) {
	goos := runtime.GOOS
	return &goos, nil
}
