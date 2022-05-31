//go:build darwin

package cicada

// RecognizeOs identifies the environment,
// as an endoflife.date product name.
func RecognizeOs() (*string, error) {
	macOs := "macos"
	return &macOs, nil
}
