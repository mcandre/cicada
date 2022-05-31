package cicada_test

import (
	"github.com/mcandre/cicada"
	"gopkg.in/yaml.v3"

	"reflect"
	"regexp"
	"testing"
)

func TestVersionQueryYAMLCodec(t *testing.T) {
	query := cicada.VersionQuery{
		Command: []string{"lsb_release", "-r"},
		Pattern: regexp.MustCompile("^Release:\\s+(?P<Version>[0-9\\.]+)$"),
	}

	queryYAML, err := yaml.Marshal(query)

	if err != nil {
		t.Fatal(err)
	}

	var query2 cicada.VersionQuery
	if err := yaml.Unmarshal(queryYAML, &query2); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(query2, query) {
		t.Errorf("Expected decoded query2: %v to equal original query: %v", query2, query)
	}
}
