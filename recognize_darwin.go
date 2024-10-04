//go:build darwin

package cicada

// EnvironmentIsLinux checks whether the current platform is Linux.
var EnvironmentIsLinux bool

// RecognizeOs identifies the environment,
// as an endoflife.date product name.
func RecognizeOs() (*string, error) {
	macOs := "macos"
	return &macOs, nil
}
