//go:build !darwin && !linux
package cicada

import (
	"runtime"
)

// RecognizeOs identifies the environment,
// as an endoflife.date product name.
func RecognizeOs() (*string, error) {
	goos := runtime.GOOS
	return &goos, nil
}
