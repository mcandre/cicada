package cicada

import (
	"github.com/MasterMinds/semver"
	"gopkg.in/yaml.v2"

	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"time"
)

// IndexUrl denotes the location of the LTS index resource.
const IndexUrl = "https://raw.githubusercontent.com/mcandre/cicada/main/cicada.yaml"

// EndOfLifeBaseUrl denotes the base location of the endoflife.date service.
const EndOfLifeBaseUrl = "https://endoflife.date/api"

// ProductsListResourceBase denotes the location of the products list resource.
const ProductsListResourceBase = "all.json"

// IndexCacheRoot denotes the cicada metadata directory base path,
// relative to the home directory.
const IndexCacheRoot = ".cicada"

// IndexCacheBase denotes the base path of the cached LTS index,
// relative to IndexCacheRoot.
const IndexCacheBase = "cicada.yaml"

// IndexProductsListBase denotes the base path of the products list file,
// relative to IndexCacheRoot.
const IndexProductsListBase = "products.json"

// IndexProductsDirBase denotes the base path of the products directory,
// relative to IndexCacheRoot.
const IndexProductsDirBase = "products"

// Index models a catalog of LTS schedules.
type Index struct {
	Debug bool `yaml:"debug,omitempty"`

	// VersionQueries denotes command line queries for retrieving component versions, in exec-like format,
	// keyed on executable base path.
	VersionQueries map[string]VersionQuery `yaml:"version_queries"`

	// components denotes version schedules,
	// keyed on component name.
	components map[string][]Schedule `yaml:"-"`
}

// IndexCacheDirPath yields the location of cicada metadata directory.
func IndexCacheDirPath() (*string, error) {
	user, err := user.Current()

	if err != nil {
		return nil, err
	}

	pth := path.Join(user.HomeDir, IndexCacheRoot)
	return &pth, nil
}

// IndexCacheConfigPath yields the location of the cicada configuration.
func IndexCacheConfigPath() (*string, error) {
	user, err := user.Current()

	if err != nil {
		return nil, err
	}

	pth := path.Join(user.HomeDir, IndexCacheRoot, IndexCacheBase)
	return &pth, nil
}

// CacheLifetimeData ensures a local copy of endoflife.date records.
func CacheLifetimeData(indexProductsListFilePath string, indexProductsDirPath string) error {
	log.Printf("Caching new product data...\n")

	fProductList, err := os.Create(indexProductsListFilePath)

	u := fmt.Sprintf("%v/%v", EndOfLifeBaseUrl, ProductsListResourceBase)
	res, err := http.Get(u)

	if err != nil {
		return err
	}

	statusCode := res.StatusCode

	if statusCode < 200 || statusCode > 299 {
		return fmt.Errorf("get: %v returned status code: %v", u, statusCode)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	defer func() {
		if err2 := fProductList.Close(); err2 != nil {
			fmt.Fprintf(os.Stderr, err2.Error())
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

		u2 := fmt.Sprintf("%v/%v", EndOfLifeBaseUrl, productBase)
		res, err2 := http.Get(u2)

		if err2 != nil {
			return err2
		}

		statusCode := res.StatusCode

		if statusCode < 200 || statusCode > 299 {
			return fmt.Errorf("get: %v returned status code: %v", u2, statusCode)
		}

		defer res.Body.Close()

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

// CacheIndex populates a cicada index.
func CacheIndex(indexDirPath string, indexCacheConfigPath string, indexProductsListFilePath string, indexProductsDirPath string) error {
	log.Printf("Caching new default configuration...\n")

	f, err := os.Create(indexCacheConfigPath)

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

	return CacheLifetimeData(indexProductsListFilePath, indexProductsDirPath)
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
	indexDirPathP, err := IndexCacheDirPath()

	if err != nil {
		return nil, err
	}

	indexDirPath := *indexDirPathP

	if err := os.MkdirAll(indexDirPath, os.ModePerm); err != nil {
		return nil, err
	}

	indexCacheConfigPathP, err := IndexCacheConfigPath()

	if err != nil {
		return nil, err
	}

	indexCacheConfigPath := *indexCacheConfigPathP

	indexProductsListFilePath := path.Join(indexDirPath, IndexProductsListBase)
	indexProductsDirPath := path.Join(indexDirPath, IndexProductsDirBase)

	_, err = os.Stat(indexCacheConfigPath)

	if update || os.IsNotExist(err) {
		if err2 := CacheIndex(indexDirPath, indexCacheConfigPath, indexProductsListFilePath, indexProductsDirPath); err2 != nil {
			return nil, err2
		}
	}

	index := new(Index)
	index.components = make(map[string][]Schedule)

	contentYAML, err := ioutil.ReadFile(indexCacheConfigPath)

	if err != nil {
		return nil, err
	}

	if err2 := yaml.UnmarshalStrict(contentYAML, index); err2 != nil {
		return nil, err2
	}

	if err2 := index.Validate(); err2 != nil {
		return nil, err2
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
		return nil, fmt.Errorf("no support schedule found for os: %v\n", identityOs)
	}

	query, ok := o.VersionQueries[identityOs]

	if !ok {
		return nil, fmt.Errorf("no version query command found for os: %v\n", identityOs)
	}

	versionString, err := query.Execute()

	if err != nil {
		return nil, err
	}

	versionP, err := semver.NewVersion(*versionString)

	if err != nil {
		return nil, err
	}

	version := *versionP

	if o.Debug {
		log.Printf("detected os %v v%v\n", identityOs, version.String())
	}

	return ScanComponent(identityOs, version, schedules, t), nil
}

func (o Index) ScanApplication(executable string, schedules []Schedule, t time.Time) (*string, error) {
	if _, err := exec.LookPath(executable); err != nil {
		return nil, nil
	}

	query, ok := o.VersionQueries[executable]

	if !ok {
		fmt.Fprintf(os.Stderr, "no version query command found for executable: %v\n", executable)
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
		log.Printf("detected application %v v%v\n", executable, version.String())
	}

	return ScanComponent(executable, *version, schedules, t), nil
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
