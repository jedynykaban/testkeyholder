package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jedynykaban/testkeyholder/model"
	"github.com/jinzhu/now"

	log "github.com/Sirupsen/logrus"
)

// Kojo allows one to deal with mysterious mitem's structure.
// Here's the deal tell me what you want to extract from the mitem and I will do it.
type Kojo interface {
	GetMitemTiniest(data json.RawMessage) (*model.MitemTiniest, error)
	GetSourceURL(data json.RawMessage) (string, error)
	GetCreationDate(data json.RawMessage) (time.Time, error)
	ConvertCreationDate(mt *model.MitemTiniest) (time.Time, error)
	GetCategory(data json.RawMessage) (string, error)
	GetCategoryPath(data json.RawMessage) (string, error)
	MakeCategoryPath(cat *model.CategoryTiniest) string
	GetAuthors(data json.RawMessage) ([]string, error)
	GetLogoURL(data json.RawMessage) (string, error)
	GetStatus(data json.RawMessage) (int, error)
	GetBody(data json.RawMessage) ([]json.RawMessage, error)
	Validate(data json.RawMessage) []error
	Process(input json.RawMessage) (json.RawMessage, error)
}

// kojoService implements Kojo interface
type kojoService struct {
}

var _ Kojo = &kojoService{}

func init() {
	// we support all the formats from the time package
	var supportedTimeFormats = []string{
		time.ANSIC,                          // "Mon Jan _2 15:04:05 2006"
		time.UnixDate,                       // "Mon Jan _2 15:04:05 MST 2006"
		time.RubyDate,                       // "Mon Jan 02 15:04:05 -0700 2006"
		time.RFC822,                         // "02 Jan 06 15:04 MST"
		time.RFC822Z,                        // "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
		time.RFC850,                         // "Monday, 02-Jan-06 15:04:05 MST"
		time.RFC1123,                        // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC1123Z,                       // "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
		time.RFC3339,                        // "2006-01-02T15:04:05Z07:00"
		time.RFC3339Nano,                    // "2006-01-02T15:04:05.999999999Z07:00"
		time.Kitchen,                        // "3:04PM"
		time.Stamp,                          // "Jan _2 15:04:05"
		time.StampMilli,                     // "Jan _2 15:04:05.000"
		time.StampMicro,                     // "Jan _2 15:04:05.000000"
		time.StampNano,                      // "Jan _2 15:04:05.000000000"
		"Mon Jan 02 2006 15:04:05 MST-0700", //Custom format for svt.se feed
	}
	for _, format := range supportedTimeFormats {
		now.TimeFormats = append(now.TimeFormats, format)
	}
}

// New - ctor like function - creates an instance of kojoService object
func NewKojo() Kojo {
	//gaService, err := ga.New("", "", "")
	// if err != nil {
	// 	log.WithError(err).Errorln("Error while fetching data from GA service")
	// }
	return &kojoService{}
}

// GetSourceURL: extracts sourcURL field from the mitem structure
func (ks *kojoService) GetSourceURL(data json.RawMessage) (string, error) {
	var mt model.MitemTiniest
	err := json.Unmarshal(data, &mt)
	if err != nil {
		return "", fmt.Errorf("Unable to unmarshal passed mitem to mitemTines, error = %s", err.Error())
	}
	if len(mt.SourceURL) == 0 {
		return "", errors.New("Either sourceURL field not present in the mitem or it is empty.")
	}
	return mt.SourceURL, nil
}

func (ks *kojoService) GetBody(data json.RawMessage) ([]json.RawMessage, error) {
	type tsukijiMitem struct {
		Body []json.RawMessage `json:"body"`
	}

	mitem := tsukijiMitem{}
	err := json.Unmarshal(data, &mitem)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Unable to decode passed mitem")
		return nil, err
	}
	return mitem.Body, nil
}

// GetAuthors: extracts authors names from the mitem structure
func (ks *kojoService) GetAuthors(data json.RawMessage) ([]string, error) {
	var rawMitem map[string]interface{}
	err := json.Unmarshal(data, &rawMitem)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Unable to decode passed mitem")
		return nil, err
	}
	// extract authors collection from the mitem as interface
	ai, ok := rawMitem["authors"]
	if !ok {
		return nil, errors.New("Unable to extract authors from the mitem. No such field found")
	}
	// convert to collection
	authors, ok := ai.([]interface{})
	if !ok {
		return nil, errors.New("Type assertion error unable to cast authors value to []interface{} type")
	}
	if len(authors) == 0 {
		return nil, nil
	}
	var ret []string
	for _, ai := range authors {
		author, ok := ai.(map[string]interface{})
		if !ok {
			return nil, errors.New("Type assertion error unable to cast author interface value to map[string]interface{} type")
		}
		// extract name field
		ni, ok := author["name"]
		if !ok {
			return nil, errors.New("Unable to extract name from the author. No such field found")
		}
		name, ok := ni.(string)
		if !ok {
			return nil, errors.New("Type assertion error unable to cast author name value to string type")
		}
		if len(name) > 0 {
			ret = append(ret, name)
		}
	}
	return ret, nil
}

// GetCreationDate: extracts date field from the mitem structure
func (ks *kojoService) GetCreationDate(data json.RawMessage) (time.Time, error) {
	var mt model.MitemTiniest
	var ret time.Time
	err := json.Unmarshal(data, &mt)
	if err != nil {
		return ret, errors.New("Unable to unmarshal passed mitem")
	}
	return ks.ConvertCreationDate(&mt)
}

// ConvertCreationDate parses date string and converts to time.Time structure
func (ks *kojoService) ConvertCreationDate(mt *model.MitemTiniest) (time.Time, error) {
	return now.Parse(mt.Date)
}

// GetCategory: extracts category field from the mitem structure
func (ks *kojoService) GetCategory(data json.RawMessage) (string, error) {
	var mt model.MitemTiniest
	err := json.Unmarshal(data, &mt)
	if err != nil {
		return "", errors.New("Unable to unmarshal passed mitem")
	}
	return mt.Category.Tier1, nil
}

// GetCategoryPath: extracts category tier1 and tier2 fields from the mitem structure, and creates full path
func (ks *kojoService) GetCategoryPath(data json.RawMessage) (string, error) {
	var mt model.MitemTiniest
	err := json.Unmarshal(data, &mt)
	if err != nil {
		return "", errors.New("Unable to unmarshal passed mitem")
	}
	return ks.MakeCategoryPath(&mt.Category), nil
}

// MakeCategoryPath: creates category path from Category structure
func (ks *kojoService) MakeCategoryPath(cat *model.CategoryTiniest) string {
	ret := cat.Tier1
	if len(ret) > 0 && len(cat.Tier2) > 0 {
		ret += ">" + cat.Tier2
	}
	return ret
}

// GetLogoURL: extracts logo URL field from the mitem structure
func (ks *kojoService) GetLogoURL(data json.RawMessage) (string, error) {
	var mt model.MitemTiniest
	err := json.Unmarshal(data, &mt)
	if err != nil {
		return "", errors.New("Unable to unmarshal passed mitem")
	}
	return mt.Meta.LogoURL, nil
}

// GetStatus: extracts status field from the mitem structure
func (ks *kojoService) GetStatus(data json.RawMessage) (int, error) {
	var mt model.MitemTiniest
	err := json.Unmarshal(data, &mt)
	if err != nil {
		return -1, errors.New("Unable to unmarshal passed mitem")
	}
	return mt.Status, nil
}

// GetMitemTiniest: converts raw data into structure
func (ks *kojoService) GetMitemTiniest(data json.RawMessage) (*model.MitemTiniest, error) {
	var mt model.MitemTiniest
	err := json.Unmarshal(data, &mt)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal passed mitem to mitemTines, error = %s", err.Error())
	}
	return &mt, nil
}

// ProcessFunc is a definition of function used to process raw mitem data
type processFunc func(data json.RawMessage) (json.RawMessage, error)

// Process calls all process functions passed as arguments
func (ks *kojoService) Process(input json.RawMessage) (json.RawMessage, error) {
	return chainProcess(input)
}

// Calls all required processing functions in chain
func chainProcess(input json.RawMessage, steps ...processFunc) (json.RawMessage, error) {
	processed := input
	var err error
	for _, step := range steps {
		processed, err = step(processed)
		if err != nil {
			return nil, err
		}
	}
	return processed, nil
}

// Validate mitems structure.
// Note we don't immediately stop on first error.
// Thus you can expect multiple error messages in the output.
func (ks *kojoService) Validate(data json.RawMessage) []error {
	log.Debug("Validating the mitem")
	var ret []error
	if len(data) <= 0 {
		ret = append(ret, errors.New("An empty mitem passed in"))
	} else {
		var mt model.MitemTiniest
		err := json.Unmarshal(data, &mt)
		if err != nil {
			ret = append(ret, errors.New("Unable to unmarshal passed mitem"))
		} else {
			ret = append(ret, mt.Validate()...)
		}
	}
	return ret
}
