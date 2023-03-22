package cicada

import (
	"github.com/MasterMinds/semver"
	"gopkg.in/yaml.v3"

	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"
)

// IndexURL denotes the location of the LTS index resource.
const IndexURL = "https://raw.githubusercontent.com/mcandre/cicada/main/cicada.yaml"

// EndOfLifeBaseURL denotes the base location of the endoflife.date service.
const EndOfLifeBaseURL = "https://endoflife.date/api"

// ProductsListResourceBase denotes the location of the products list resource.
const ProductsListResourceBase = "all.json"

// IndexCacheRoot denotes the cicada metadata directory base path,
// relative to the home directory.
const IndexCacheRoot = ".cicada"

// IndexCacheBase denotes the base path of the cached LTS index,
// relative to the current working directory.
//
// For example, a software project top level directory.
const IndexCacheBase = "cicada.yaml"

// IndexProductsListBase denotes the base path of the products list file,
// relative to IndexCacheRoot.
const IndexProductsListBase = "products.json"

// IndexProductsDirBase denotes the base path of the products directory,
// relative to IndexCacheRoot.
const IndexProductsDirBase = "products"

// Index models a catalog of LTS schedules.
type Index struct {
	// Debug enables additional logging (default: false).
	Debug bool `yaml:"debug,omitempty"`

	// Quiet skips system executables (default: false).
	Quiet bool `yaml:"quiet,omitempty"`

	// LeadMonths provides a margin of time to migrate
	// before a support timeline formally ends.
	//
	// Not too short that developers fail to migrate,
	// not too long that developers forget to migrate.
	//
	// Negative values are treated as a reset to default value.
	//
	// (default: 1)
	LeadMonths int `yaml:"lead_months,omitempty"`

	// VersionQueries denotes command line queries for retrieving component versions, in exec-like format,
	// keyed on executable base path.
	VersionQueries map[string]VersionQuery `yaml:"version_queries"`

	// components denotes version schedules,
	// keyed on component name.
	components map[string][]Schedule `yaml:"-"`
}

// IndexCacheDirPath yields the location of cicada metadata directory.
func IndexCacheDirPath(cwd string) (*string, error) {
	pth := path.Join(cwd, IndexCacheRoot)
	return &pth, nil
}

// IndexCacheConfigPath yields the location of the cicada configuration.
func IndexCacheConfigPath(cwd string) (*string, error) {
	pth := path.Join(cwd, IndexCacheBase)
	return &pth, nil
}

// CacheLifetimeData ensures a local copy of endoflife.date records.
func CacheLifetimeData(indexProductsListFilePath string, indexProductsDirPath string) error {
	log.Println("Caching new product data...")

	fProductList, err := os.Create(indexProductsListFilePath)

	if err != nil {
		return err
	}

	u := fmt.Sprintf("%v/%v", EndOfLifeBaseURL, ProductsListResourceBase)
	res, err := http.Get(u)

	if err != nil {
		return err
	}

	statusCode := res.StatusCode

	if statusCode < 200 || statusCode > 299 {
		return fmt.Errorf("get: %v returned status code: %v", u, statusCode)
	}

	defer func() {
		if err2 := res.Body.Close(); err2 != nil {
			log.Print(err2)
		}
	}()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	defer func() {
		if err2 := fProductList.Close(); err2 != nil {
			log.Print(err2)
		}
	}()

	if _, err2 := fProductList.Write(body); err2 != nil {
		return err2
	}

	var products []string
	if err2 := json.Unmarshal(body, &products); err2 != nil {
		return err2
	}

	if err2 := os.MkdirAll(indexProductsDirPath, os.ModePerm); err2 != nil {
		return err2
	}

	for _, product := range products {
		productBase := fmt.Sprintf("%v.json", product)
		productFilePath := path.Join(indexProductsDirPath, productBase)
		fProductDetail, err2 := os.Create(productFilePath)

		if err2 != nil {
			return err2
		}

		u2 := fmt.Sprintf("%v/%v", EndOfLifeBaseURL, productBase)
		res, err2 := http.Get(u2)

		if err2 != nil {
			return err2
		}

		statusCode := res.StatusCode

		if statusCode < 200 || statusCode > 299 {
			return fmt.Errorf("get: %v returned status code: %v", u2, statusCode)
		}

		defer func() {
			if err3 := res.Body.Close(); err3 != nil {
				log.Print(err3)
			}
		}()

		body, err2 := io.ReadAll(res.Body)

		if err2 != nil {
			return err2
		}

		if _, err3 := fProductDetail.Write(body); err3 != nil {
			return err3
		}
	}

	return nil
}

// ValidateVersionQueries ensures version query data integrity.
func (o Index) ValidateVersionQueries() error {
	for component, query := range o.VersionQueries {
		if len(query.Command) == 0 {
			return fmt.Errorf("%v has an empty version query", component)
		}
	}

	return nil
}

// Validate ensures data integrity.
func (o Index) Validate() error {
	return o.ValidateVersionQueries()
}

// Load generates a partial LTS index.
func Load(update bool) (*Index, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	indexDirPathP, err := IndexCacheDirPath(cwd)

	if err != nil {
		return nil, err
	}

	indexDirPath := *indexDirPathP

	if err2 := os.MkdirAll(indexDirPath, os.ModePerm); err2 != nil {
		return nil, err2
	}

	indexCacheConfigPathP, err := IndexCacheConfigPath(cwd)

	if err != nil {
		return nil, err
	}

	indexCacheConfigPath := *indexCacheConfigPathP

	if _, err2 := os.Stat(indexCacheConfigPath); os.IsNotExist(err2) {
		return nil, fmt.Errorf("missing configuration: %v", IndexCacheBase)
	}

	indexProductsListFilePath := path.Join(indexDirPath, IndexProductsListBase)
	indexProductsDirPath := path.Join(indexDirPath, IndexProductsDirBase)

	_, err = os.Stat(indexProductsListFilePath)

	if update || os.IsNotExist(err) {
		if err2 := CacheLifetimeData(indexProductsListFilePath, indexProductsDirPath); err2 != nil {
			return nil, err2
		}
	}

	index := new(Index)
	index.components = make(map[string][]Schedule)

	contentYAML, err := os.ReadFile(indexCacheConfigPath)

	if err != nil {
		return nil, err
	}

	if err2 := yaml.Unmarshal(contentYAML, index); err2 != nil {
		return nil, err2
	}

	if err2 := index.Validate(); err2 != nil {
		return nil, err2
	}

	if index.LeadMonths < 0 {
		index.LeadMonths = DefaultLeadMonths
	}

	productListBuf, err := os.ReadFile(indexProductsListFilePath)

	if err != nil {
		return nil, err
	}

	var products []string
	if err2 := json.Unmarshal(productListBuf, &products); err2 != nil {
		return nil, err2
	}

	for _, product := range products {
		productDetailPath := fmt.Sprintf("%v.json", path.Join(indexProductsDirPath, product))
		productDetailBuf, err2 := os.ReadFile(productDetailPath)

		if err2 != nil {
			return nil, err2
		}

		var records ProductRecords
		if err2 := json.Unmarshal(productDetailBuf, &records); err2 != nil {
			return nil, err2
		}

		schedules, err := ProductRecordsToSchedules(product, records)

		if err != nil {
			return nil, err
		}

		index.components[product] = schedules
	}

	return index, nil
}

// ScanOs analyzes operating system for any LTS concerns.
func (o Index) ScanOs(t time.Time) (*string, error) {
	identityOsP, err := RecognizeOs()

	if err != nil {
		return nil, err
	}

	identityOs := *identityOsP
	schedules, ok := o.components[identityOs]

	if !ok {
		log.Printf("no known support schedule found for os: %v", identityOs)
		return nil, nil
	}

	query, ok := o.VersionQueries[identityOs]

	if !ok {
		log.Printf("no known version query command found for os: %v", identityOs)
		return nil, nil
	}

	versionString, err := query.Execute()

	if err != nil {
		return nil, err
	}

	if versionString == nil {
		log.Fatalf("unable to identify version for os: %v", identityOs)
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

// ScanApplication checks executables for non-LTS versions.
//
// If the executable is not found, a warning may not be emitted.
func (o Index) ScanApplication(app string, schedules []Schedule, t time.Time) (*string, error) {
	if IsOperatingSystem(app) {
		return nil, nil
	}

	query, ok := o.VersionQueries[app]

	if !ok {
		if o.Debug {
			log.Printf("no version query command found for application: %v\n", app)
		}

		return nil, nil
	}

	executable := query.Command[0]

	executablePath, err := exec.LookPath(executable)

	if err != nil {
		if o.Debug {
			log.Printf("executable not found: %v for application %v; skipping\n", executable, app)
		}

		return nil, nil
	}

	if o.Quiet && IsSystemExecutable(executablePath) {
		if o.Debug {
			log.Printf("executable: %v found in system path; skipping\n", executable)
		}

		return nil, nil
	}

	versionString, err := query.Execute()

	if err != nil {
		return nil, err
	}

	if versionString == nil {
		return nil, nil
	}

	version, err := semver.NewVersion(*versionString)

	if err != nil {
		return nil, err
	}

	if o.Debug {
		log.Printf("detected application: %v v%v\n", app, version.String())
	}

	return ScanComponent(app, *version, schedules, t), nil
}

// ScanApplications analyzes applications for any LTS concerns.
func (o Index) ScanApplications(t time.Time) ([]string, error) {
	var warnings []string

	for executable, schedules := range o.components {
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
	t := tNow.AddDate(0, o.LeadMonths, 0)
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

// Clean removes artifacts created during cicada runs.
func Clean() error {
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	indexCacheDirPath, err := IndexCacheDirPath(cwd)

	if err != nil {
		return err
	}

	return os.RemoveAll(*indexCacheDirPath)
}
