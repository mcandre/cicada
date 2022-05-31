package cicada_test

import (
	"github.com/mcandre/cicada"
	"gopkg.in/yaml.v3"

	"reflect"
	"regexp"
	"testing"
)

func TestIndexYAMLCodec(t *testing.T) {
	index := cicada.Index{
		Debug: true,
		VersionQueries: map[string]cicada.VersionQuery{
			"ruby": cicada.VersionQuery{
				Command: []string{"lsb_release", "-r"},
				Pattern: regexp.MustCompile("^Release:\\s+(?P<Version>[0-9\\.]+)$"),
			},
		},
	}

	indexYAML, err := yaml.Marshal(index)

	if err != nil {
		t.Fatal(err)
	}

	var index2 cicada.Index
	if err := yaml.Unmarshal(indexYAML, &index2); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(index2, index) {
		t.Errorf("Expected decoded index2: %v to equal original query: %v", index2, index)
	}
}
