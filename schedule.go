package cicada

import (
	"github.com/MasterMinds/semver"
	"gopkg.in/yaml.v2"

	"fmt"
	"time"
)

// RFC3339DateFormat presents a prefix of the RFC3339 timestamp format.
const RFC3339DateFormat = "2006-01-02"

// Schedule models LTS series.
type Schedule struct {
	// Name denotes a software component:
	// Either a GOOS value or an executable base path.
	Name string `yaml:"name"`

	// Version denotes a software release series.
	Version semver.Version `yaml:"version"`

	// Expiration denotes a termination timestamp.
	Expiration *time.Time `yaml:"expiration,omitempty"`
}

// Match reports whether a schedule applies to the given software component version.
func (o Schedule) Match(v semver.Version) bool {
	if v.Major() != o.Version.Major() {
		return false
	}

	minor := o.Version.Minor()

	if minor == 0 {
		return true
	}

	return v.Minor() == minor
}

// MarshalYAML encodes schedules.
func (o Schedule) MarshalYAML() ([]byte, error) {
	type ScheduleAlias struct {
		Name       string `yaml:"name"`
		Version    string `yaml:"version"`
		Expiration string `yaml:"expiration,omitempty"`
	}

	var aux ScheduleAlias
	aux.Name = o.Name
	aux.Version = o.Version.String()

	if o.Expiration != nil {
		aux.Expiration = o.Expiration.Format(RFC3339DateFormat)
	}

	return yaml.Marshal(aux)
}

// UnmarshalYAML decodes schedules.
func (o *Schedule) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ScheduleAlias struct {
		Name       string `yaml:"name"`
		Version    string `yaml:"version"`
		Expiration string `yaml:"expiration,omitempty"`
	}

	var aux ScheduleAlias

	if err := unmarshal(&aux); err != nil {
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
func ScanComponent(name string, version semver.Version, schedules []Schedule, t time.Time) *string {
	for _, schedule := range schedules {
		if !schedule.Match(version) {
			continue
		}

		if schedule.Expiration != nil {
			expiration := *schedule.Expiration

			if t.After(expiration) {
				message := fmt.Sprintf("end of life for %v v%v on %v", name, version.String(), expiration.Format(RFC3339DateFormat))
				return &message
			}
		}
	}

	return nil
}
