package cicada

import (
	"gopkg.in/yaml.v3"

	"bufio"
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
	outputBytes, err := cmd.Output()

	if err != nil {
		return nil, nil
	}

	outputString := string(outputBytes)
	versionString := strings.TrimRight(outputString, "\r\n")

	if o.Pattern != nil {
		scanner := bufio.NewScanner(strings.NewReader(outputString))

		var foundVersion bool

		for scanner.Scan() {
			line := scanner.Text()
			match := o.Pattern.FindStringSubmatch(line)
			versionIndex := o.Pattern.SubexpIndex("Version")

			if len(match) > versionIndex {
				foundVersion = true
				versionString = match[versionIndex]
				break
			}
		}

		if !foundVersion {
			return nil, nil
		}
	}

	return &versionString, nil
}
