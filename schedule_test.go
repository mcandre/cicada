package cicada_test

import (
	"github.com/Masterminds/semver"
	"github.com/mcandre/cicada"
	"gopkg.in/yaml.v3"

	"reflect"
	"testing"
	"time"
)

func TestScheduleYAMLCodec(t *testing.T) {
	version, err := semver.NewVersion("2.6")

	if err != nil {
		t.Fatal(err)
	}

	exp, err := time.Parse(cicada.RFC3339DateFormat, "2022-03-31")

	if err != nil {
		t.Fatal(err)
	}

	schedule := cicada.Schedule{
		Name:       "Ruby",
		Version:    *version,
		Expiration: &exp,
	}

	scheduleYAML, err := yaml.Marshal(schedule)

	if err != nil {
		t.Fatal(err)
	}

	var schedule2 cicada.Schedule
	if err := yaml.Unmarshal(scheduleYAML, &schedule2); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(schedule2, schedule) {
		t.Errorf("Expected decoded schedule2: %v to equal original schedule: %v", schedule2, schedule)
	}
}
