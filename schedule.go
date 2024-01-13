package cicada

import (
	"github.com/MasterMinds/semver"
	"gopkg.in/yaml.v3"

	"fmt"
	"regexp"
	"strings"
	"time"
)

// RFC3339DateFormat presents a prefix of the RFC3339 timestamp format.
const RFC3339DateFormat = "2006-01-02"

// Schedule models LTS series.
type Schedule struct {
	// Name denotes a software component:
	// Either a GOOS value or an executable base path.
	Name string `json:"name" yaml:"name"`

	// Codename denotes a version nickname.
	Codename string `json:"codename" yaml:"codename"`

	// Version denotes a software release series.
	// Only the major and minor are included in end of life calculations.
	// Zero minor is treated as matching any minor.
	Version semver.Version `json:"version" yaml:"version"`

	// Expiration denotes a termination timestamp.
	//
	// nil indicates no known expiration.
	//
	// (default: nil)
	Expiration *time.Time `json:"expiration,omitempty" yaml:"expiration,omitempty"`
}

// Match reports whether a schedule applies to the given software component version or codename.
//
// specificity indicates the number of elements in the original v string.
//
// For example, original version string "1" has specificity 1.
// Original version string "1.1" has specificity 2.
// Original version string "1.1.1" has specificity 3.
// And so on.
//
// Note that degenerate versions may not necessarily behave as expected.
// For example, ".1" (corresponding with "0.1"),
// Or "1." (corresponding with "1.0").
func (o Schedule) Match(version *semver.Version, specificity int, codename string) bool {
	if codename != "" {
		return regexp.MustCompile(fmt.Sprintf("(?i)%s", codename)).MatchString(o.Codename)
	}

	if version == nil {
		return true
	}

	v := *version

	if v.Major() != o.Version.Major() {
		return false
	}

	if specificity < 1 {
		return true
	}

	return v.Minor() == o.Version.Minor()
}

// MarshalYAML encodes schedules.
func (o Schedule) MarshalYAML() (interface{}, error) {
	type ScheduleAlias struct {
		Name       string `yaml:"name"`
		Version    string `yaml:"version"`
		Expiration string `yaml:"expiration,omitempty"`
	}

	var aux ScheduleAlias
	aux.Name = o.Name
	aux.Version = o.Version.Original()

	if o.Expiration != nil {
		aux.Expiration = o.Expiration.Format(RFC3339DateFormat)
	}

	return aux, nil
}

// UnmarshalYAML decodes schedules.
func (o *Schedule) UnmarshalYAML(value *yaml.Node) error {
	type ScheduleAlias struct {
		Name       string `json:"name" yaml:"name"`
		Version    string `json:"version" yaml:"version"`
		Expiration string `json:"expiration,omitempty" yaml:"expiration,omitempty"`
	}

	var aux ScheduleAlias

	if err := value.Decode(&aux); err != nil {
		return err
	}

	if aux.Expiration != "" {
		t, err := time.Parse(RFC3339DateFormat, aux.Expiration)

		if err != nil {
			return err
		}

		o.Expiration = &t
	}

	o.Name = aux.Name
	version, err := semver.NewVersion(aux.Version)

	if err != nil {
		return err
	}

	o.Version = *version
	return nil
}

// ScanComponent checks whether the given component is end of life.
func ScanComponent(name string, version *semver.Version, codename string, schedules []Schedule, t time.Time) *string {
	var specificity int

	if version != nil {
		specificity = strings.Count(version.Original(), ".")
	}

	for _, schedule := range schedules {
		if !schedule.Match(version, specificity, codename) {
			continue
		}

		if schedule.Expiration != nil {
			expiration := *schedule.Expiration

			versionString := codename

			if version != nil {
				versionString = version.String()
			}

			if t.Equal(expiration) || t.After(expiration) {
				message := fmt.Sprintf("end of life for %v %v on %v", name, versionString, expiration.Format(RFC3339DateFormat))
				return &message
			}
		}
	}

	return nil
}
