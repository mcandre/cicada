package cicada

import (
	"github.com/MasterMinds/semver"
	"gopkg.in/yaml.v2"

	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path"
	"runtime"
	"time"
)

// IndexUrl denotes the location of the LTS index resource.
const IndexUrl = "https://raw.githubusercontent.com/mcandre/cicada/main/cicada.yaml"

// IndexCacheBase denotes the base path of the cached LTS index.
const IndexCacheBase = "cicada.yaml"

// Index models a catalog of LTS schedules.
type Index struct {
	// OperatingSystems denotes operating system schedules.
	OperatingSystems map[string][]Schedule `yaml:"operating_systems"`

	// ProgrammingLanguages denotes programming langauge schedules.
	ProgrammingLanguages map[string][]Schedule `yaml:"programming_languages"`
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

	return index, nil
}

// ScanOs analyzes operating system for any LTS concerns.
func (o Index) ScanOs(t time.Time) (*string, error) {
	identityOs := runtime.GOOS
	schedules, ok := o.OperatingSystems[identityOs]

	if !ok {
		return nil, fmt.Errorf("no support schedule found for os: %v\n", identityOs)
	}

	versionString, err := GetOsVersion()

	if err != nil {
		return nil, err
	}

	version, err := semver.NewVersion(*versionString)

	if err != nil {
		return nil, err
	}

	versionMajor := version.Major()
	versionMinor := version.Minor()

	var foundSchedule bool

	for _, schedule := range schedules {
		scheduleVersion, err2 := semver.NewVersion(schedule.Version)

		if err2 != nil {
			return nil, err2
		}

		scheduleVersionMajor := scheduleVersion.Major()

		if versionMajor != scheduleVersionMajor {
			continue
		}

		scheduleVersionMinor := scheduleVersion.Minor()

		if scheduleVersionMinor != 0 && versionMinor != scheduleVersionMinor {
			continue
		}

		foundSchedule = true

		if schedule.Expiration != nil {
			expiration := *schedule.Expiration

			if t.After(expiration) {
				message := fmt.Sprintf("end of life for os %v %v on %v", identityOs, versionString, expiration)
				return &message, nil
			}
		}
	}

	if !foundSchedule {
		return nil, fmt.Errorf("no matching support schedule found for os %v version %v\n", identityOs, versionString)
	}

	return nil, nil
}

// ScanProgrammingLanguages analyzes programming languages for any LTS concerns.
func (o Index) ScanProgrammingLanguages(t time.Time) ([]string, error) {
	// ...
	return nil, nil
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

	resultsProgrammingLanguages, err := o.ScanProgrammingLanguages(t)

	if err != nil {
		return nil, err
	}

	warnings = append(warnings, resultsProgrammingLanguages...)
	return warnings, nil
}
