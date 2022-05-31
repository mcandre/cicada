package cicada

import (
	"gopkg.in/yaml.v3"

	"os/exec"
	"regexp"
	"strings"
)

// VersionQuery models commands for extracting software component version information.
type VersionQuery struct {
	// Command denotes an exec-like command line instruction.
	//
	// Command output is always right trimmed.
	Command []string `yaml:"command"`

	// Pattern denotes an optional expression for
	// capturing version strings within
	// larger, complex output buffers.
	//
	// nil indicates the full command output,
	// sans right trim,
	// is treated as a semver version string.
	//
	// (default: nil)
	Pattern *regexp.Regexp `yaml:"pattern,omitempty"`
}

// MarshalYAML encodes version queries.
func (o VersionQuery) MarshalYAML() (interface{}, error) {
	type VersionQueryAlias struct {
		Command []string `yaml:"command"`
		Pattern *string  `yaml:"pattern,omitempty"`
	}

	var aux VersionQueryAlias
	aux.Command = o.Command
	patternString := o.Pattern.String()
	aux.Pattern = &patternString
	return aux, nil
}

// UnmarshalYAML decodes version queries.
func (o *VersionQuery) UnmarshalYAML(value *yaml.Node) error {
	type VersionQueryAlias struct {
		Command []string `yaml:"command"`
		Pattern *string  `yaml:"pattern,omitempty"`
	}

	var aux VersionQueryAlias

	if err := value.Decode(&aux); err != nil {
		return err
	}

	if aux.Pattern != nil {
		patternString := *aux.Pattern

		if patternString != "" {
			pattern, err := regexp.Compile(patternString)

			if err != nil {
				return err
			}

			o.Pattern = pattern
		}
	}

	o.Command = aux.Command
	return nil
}

// Execute retrieves software component versions.
func (o VersionQuery) Execute() (*string, error) {
	command, args := o.Command[0], o.Command[1:]
	cmd := exec.Command(command, args...)
	versionBytes, err := cmd.Output()

	if err != nil {
		return nil, nil
	}

	versionString := string(versionBytes)
	versionString = strings.TrimRight(versionString, "\r\n")

	if o.Pattern != nil {
		matches := o.Pattern.FindStringSubmatch(versionString)
		versionIndex := o.Pattern.SubexpIndex("Version")

		if len(matches) < versionIndex+1 {
			return nil, nil
		}

		versionString = matches[versionIndex]
	}

	return &versionString, nil
}
