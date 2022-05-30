package cicada

import (
	"gopkg.in/yaml.v2"

	"time"
)

// RFC3339DateFormat presents a prefix of the RFC3339 timestamp format.
const RFC3339DateFormat = "2006-01-02"

// Schedule models LTS series.
type Schedule struct {
	// Name denotes a software component.
	Name string `yaml:"name"`

	// Version denotes a semver release range constraint.
	Version string `yaml:"version"`

	// Expiration denotes a termination timestamp.
	Expiration *time.Time `yaml:"expiration,omitempty"`
}

func (o Schedule) MarshalYAML() ([]byte, error) {
	type ScheduleAlias struct {
		Name string `yaml:"name"`
		Version string `yaml:"version"`
		Expiration string `yaml:"expiration,omitempty"`
	}

	var aux ScheduleAlias
	aux.Name = o.Name
	aux.Version = o.Version

	if o.Expiration != nil {
		aux.Expiration = o.Expiration.Format(RFC3339DateFormat)
	}

	return yaml.Marshal(aux)
}

func (o *Schedule) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ScheduleAlias struct {
		Name string `yaml:"name"`
		Version string `yaml:"version"`
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
	o.Version = aux.Version
	return nil
}
