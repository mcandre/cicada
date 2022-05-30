package cicada

import (
	"github.com/MasterMinds/semver"
	"gopkg.in/yaml.v2"

	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"runtime"
	"strings"
	"time"
)

// IndexUrl denotes the location of the LTS index resource.
const IndexUrl = "https://raw.githubusercontent.com/mcandre/cicada/main/cicada.yaml"

// IndexCacheBase denotes the base path of the cached LTS index.
const IndexCacheBase = "cicada.yaml"

// Index models a catalog of LTS schedules.
type Index struct {
	Debug bool `yaml:"debug,omitempty"`

	// OperatingSystems denotes operating system schedules.
	OperatingSystems map[string][]Schedule `yaml:"operating_systems"`

	// Applications denotes application schedules,
	// keyed on executable base path.
	Applications map[string][]Schedule `yaml:"applications"`

	// VersionQueries denotes command line queries for retrieving component versions, in exec-like format,
	// keyed on executable base path.
	VersionQueries map[string][]string `yaml:"version_queries"`
}

// IndexCachePath denotes the location of the cached LTS index.
func IndexCachePath() (*string, error) {
	user, err := user.Current()

	if err != nil {
		return nil, err
	}

	pth := path.Join(user.HomeDir, IndexCacheBase)

	return &pth, nil
}

// CacheIndex populates a local LTS index.
func CacheIndex(indexCachePath string) error {
	f, err := os.Create(indexCachePath)

	if err != nil {
		return err
	}

	res, err := http.Get(IndexUrl)

	if err != nil {
		return err
	}

	statusCode := res.StatusCode

	if statusCode < 200 || statusCode > 299 {
		return fmt.Errorf("get: %v returned status code: %v", IndexUrl, statusCode)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	defer func() {
		if err2 := f.Close(); err2 != nil {
			fmt.Fprintf(os.Stderr, err2.Error())
		}
	}()

	if _, err := f.Write(body); err != nil {
		return err
	}

	return nil
}

// ValidateVersionQueries ensures version query data integrity.
func (o Index) ValidateVersionQueries() error {
	for component, query := range o.VersionQueries {
		if len(query) == 0 {
			return fmt.Errorf("%v has an empty version query", component)
		}
	}

	return nil
}

// Validate ensures data integrity.
func (o Index) Validate() error {
	return o.ValidateVersionQueries()
}

// Load generates an LTS index.
func Load(update bool) (*Index, error) {
	indexCachePathP, err := IndexCachePath()

	if err != nil {
		return nil, err
	}

	indexCachePath := *indexCachePathP

	_, err = os.Stat(indexCachePath)

	if update || os.IsNotExist(err) {
		if err2 := CacheIndex(indexCachePath); err2 != nil {
			return nil, err2
		}
	}

	index := new(Index)

	contentYAML, err := ioutil.ReadFile(indexCachePath)

	if err != nil {
		return nil, err
	}

	if err2 := yaml.UnmarshalStrict(contentYAML, index); err2 != nil {
		return nil, err2
	}

	if err2 := index.Validate(); err2 != nil {
		return nil, err2
	}

	return index, nil
}

// QueryVersion extracts software component versions.
func QueryVersion(query []string) (*string, error) {
	command, args := query[0], query[1:]
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	versionBytes, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	versionString := string(versionBytes)
	versionString = strings.TrimRight(versionString, "\r\n")
	return &versionString, nil
}

// ScanOs analyzes operating system for any LTS concerns.
func (o Index) ScanOs(t time.Time) (*string, error) {
	identityOs := runtime.GOOS
	schedules, ok := o.OperatingSystems[identityOs]

	if !ok {
		return nil, fmt.Errorf("no support schedule found for os: %v\n", identityOs)
	}

	query, ok := o.VersionQueries[identityOs]

	if !ok {
		return nil, fmt.Errorf("no version query command found for os: %v\n", identityOs)
	}

	versionString, err := QueryVersion(query)

	if err != nil {
		return nil, err
	}

	versionP, err := semver.NewVersion(*versionString)

	if err != nil {
		return nil, err
	}

	version := *versionP

	if o.Debug {
		log.Printf("detected os: %v v%v\n", identityOs, version.String())
	}

	return ScanComponent(identityOs, version, schedules, t), nil
}

func (o Index) ScanApplication(executable string, schedules []Schedule, t time.Time) (*string, error) {
	_, err := exec.LookPath(executable)

	if err != nil {
		return nil, err
	}

	query, ok := o.VersionQueries[executable]

	if !ok {
		return nil, fmt.Errorf("no version query command found for executable: %v\n", executable)
	}

	versionString, err := QueryVersion(query)

	if err != nil {
		return nil, err
	}

	version, err := semver.NewVersion(*versionString)

	if err != nil {
		return nil, err
	}

	if o.Debug {
		log.Printf("detected application %v v%v\n", executable, version.String())
	}

	return ScanComponent(executable, *version, schedules, t), nil
}

// ScanApplications analyzes applications for any LTS concerns.
func (o Index) ScanApplications(t time.Time) ([]string, error) {
	var warnings []string

	for executable, schedules := range o.Applications {
		warning, err := o.ScanApplication(executable, schedules, t)

		if err != nil {
			return nil, err
		}

		if warning != nil {
			warnings = append(warnings, *warning)
		}
	}

	return warnings, nil
}

// Scan generates reports.
func (o Index) Scan() ([]string, error) {
	var warnings []string
	tNow := time.Now()
	t := tNow.AddDate(0, -1 * LeadMonths, 0)
	warningOs, err := o.ScanOs(t)

	if err != nil {
		return nil, err
	}

	if warningOs != nil {
		warnings = append(warnings, *warningOs)
	}

	resultsApplications, err := o.ScanApplications(t)

	if err != nil {
		return nil, err
	}

	warnings = append(warnings, resultsApplications...)
	return warnings, nil
}
