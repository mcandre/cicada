// +build darwin
package cicada

import (
	"os"
	"os/exec"
	"strings"
)

func GetOsVersion() (*string, error) {
	cmd := exec.Command("sw_vers")
	cmd.Args = []string{"sw_vers", "-productVersion"}
	cmd.Stderr = os.Stderr

	versionBytes, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	version := string(versionBytes)
	version = strings.TrimRight(version, "\r\n")
	return &version, nil
}
