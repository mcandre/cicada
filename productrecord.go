package cicada

import (
	"github.com/MasterMinds/semver"

	"regexp"
	"strconv"
	"time"
)

// SemVerPattern matches semantic versions.
var SemVerPattern = regexp.MustCompile(`^[0-9\.]+$`)

// ProductRecords models endoflife.date product detail records.
type ProductRecords []map[string]interface{}

// ProductRecordsToSchedules converts ProductRecords to Schedule arrays.
func ProductRecordsToSchedules(name string, records ProductRecords) ([]Schedule, error) {
	var schedules []Schedule

	for _, record := range records {
		cycle := record["cycle"]
		var version *semver.Version

		if c, ok := cycle.(string); ok {
			if !SemVerPattern.MatchString(c) {
				continue
			}

			v, err := semver.NewVersion(c)

			if err != nil {
				return nil, err
			}

			version = v
		} else if c, ok := cycle.(int); ok {
			majorString := strconv.Itoa(c)
			v, err := semver.NewVersion(majorString)

			if err != nil {
				return nil, err
			}

			version = v
		}

		schedule := Schedule{
			Name:    name,
			Version: *version,
		}

		eol := record["eol"]

		var expiration *time.Time

		if e, ok := eol.(string); ok {
			exp, err := time.Parse(RFC3339DateFormat, e)

			if err != nil {
				return nil, err
			}

			expiration = &exp
		}

		if expiration != nil {
			schedule.Expiration = expiration
		}

		schedules = append(schedules, schedule)
	}

	return schedules, nil
}
