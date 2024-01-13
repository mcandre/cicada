//go:build darwin

package cicada

var EnvironmentIsLinux bool = false

// RecognizeOs identifies the environment,
// as an endoflife.date product name.
func RecognizeOs() (*string, error) {
	macOs := "macos"
	return &macOs, nil
}
