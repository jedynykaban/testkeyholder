package model

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jinzhu/now"
)

const (
	bodyElementParagrahType = "paragraph"
	bodyElementH1Type       = "h1"
	bodyElementH2Type       = "h2"
	bodyElementH3Type       = "h3"
	bodyElementH4Type       = "h4"
	bodyElementH5Type       = "h5"
	bodyElementH6Type       = "h6"
	bodyElementInfoType     = "info"
	bodyElementSubheadType  = "subhead"
	bodyElementImageType    = "image"
	bodyElementVideoType    = "video"
	bodyElementGalleryType  = "gallery"

	supportedVideoTypeVimeo   = "vimeo"
	supportedVideoTypeYoutube = "youtube"
)

// MitemTiniest contains mitem metadata
type MitemTiniest struct {
	SourceURL    string            `json:"sourceURL"`
	Date         string            `json:"date"`
	Type         string            `json:"type"`
	LicenseType  string            `json:"licensetype"`
	LicenseText  string            `json:"licensetext"`
	LicensePromo string            `json:"licensepromo"`
	MainImage    bodyImageTiniest  `json:"mainimage"`
	Headline     string            `json:"headline"`
	Price        priceTiniest      `json:"price"`
	Category     CategoryTiniest   `json:"category"`
	AdsPolicy    AdsPolicyTiniest  `json:"adspolicy"`
	Meta         MetaTiniest       `json:"meta,omitempty"`
	Status       int               `json:"status,omitempty"`
	Body         []json.RawMessage `json:"body"`
}

type priceTiniest struct {
	Value    int    `json:"value"`
	Currency string `json:"currency"`
}

type CategoryTiniest struct {
	Tier1 string `json:"tier1"`
	Tier2 string `json:"tier2"`
}

type bodyElement struct {
	Type string `json:"type"`
}

type bodyCommonTiniest struct {
	bodyElement
	Content string `json:"content"`
}

type bodyImageTiniest struct {
	bodyElement
	Source  string `json:"source"`
	Caption string `json:"caption"`
	Height  int    `json:"height"`
	Width   int    `json:"width"`
}

type bodyVideoTiniest struct {
	bodyElement
	Source    string `json:"source"`
	VideoType string `json:"videoType"`
}

type bodyGalleryTiniest struct {
	bodyElement
	Body []json.RawMessage `json:"body"`
}

// AdsPolicyTiniest is the tiniest policy required
type AdsPolicyTiniest struct {
	On     bool `json:"on"`
	MaxAds int  `json:"maxAds"`
}

// MetaTiniest is the tiniest meta required
type MetaTiniest struct {
	LogoURL       string `json:"logoURL,omitempty"`
	MosaiqPrimary bool   `json:"mosaiqPrimary,omitempty"`
	UserEdited    bool   `json:"userEdited,omitempty"`
	Tags          []Tag  `json:"tags,omitempty"`
}

// Validate checks agains all mandatory fields in tiniest mitemTiniest
// and also validate it's body elements
func (m *MitemTiniest) Validate() []error {
	var ret []error
	if len(m.SourceURL) == 0 {
		ret = append(ret, errors.New("Mandatory field sourceURL is empty"))
	}
	if len(m.Date) == 0 {
		ret = append(ret, errors.New("Mandatory field date is empty"))
	} else {
		// TODO: We support two formats, add more
		if _, err := now.Parse(m.Date); err != nil {
			ret = append(ret, errors.New("Mandatory field date is in unsupported format: "+err.Error()))
		}
	}
	if len(m.Type) == 0 {
		ret = append(ret, errors.New("Mandatory field type is empty"))
	}
	if len(m.LicenseType) == 0 {
		ret = append(ret, errors.New("Mandatory field license type is empty"))
	} else {
		const e string = "editorial"
		const s string = "sponsored"
		if m.LicenseType != e && m.LicenseType != s {
			m := fmt.Sprintf("Unsupported license type got = %s, want = %s or %s", m.LicenseType, e, s)
			ret = append(ret, errors.New(m))
		}
	}
	if len(m.MainImage.Source) == 0 {
		ret = append(ret, errors.New("Mandatory field mainimage.source is empty"))
	}
	if len(m.Headline) == 0 {
		ret = append(ret, errors.New("Mandatory field headline is empty"))
	}
	if len(m.Body) == 0 {
		ret = append(ret, errors.New("Mandatory field body is empty"))
	} else {
		ret = append(ret, validateBody(m.Body)...)
	}
	return ret
}

func validateBody(datas []json.RawMessage) []error {
	var ret []error
	for _, data := range datas {
		var element bodyCommonTiniest
		err := json.Unmarshal(data, &element)
		if err != nil {
			ret = append(ret, errors.New("Unable to unmarshal body element"))
		} else {
			if len(element.Type) == 0 {
				ret = append(ret, errors.New("Mandatory field type is empty in body element"))
			} else {
				ret = append(ret, validateBodyElement(data, element.Type)...)
			}
		}
	}
	return ret
}

func validateBodyElement(data json.RawMessage, elementType string) []error {
	var ret []error
	switch elementType {
	case bodyElementParagrahType,
		bodyElementH1Type,
		bodyElementH2Type,
		bodyElementH3Type,
		bodyElementH4Type,
		bodyElementH5Type,
		bodyElementH6Type,
		bodyElementInfoType,
		bodyElementSubheadType:
		ret = append(ret, validateBodyElementCommon(data, elementType)...)

	case bodyElementImageType:
		ret = append(ret, validateBodyElementImage(data, elementType)...)

	case bodyElementVideoType:
		ret = append(ret, validateBodyElementVideo(data, elementType)...)

	case bodyElementGalleryType:
		ret = append(ret, validateBodyElementGallery(data, elementType)...)

	default:
		// ignore other element types in validation
	}
	return ret
}

func validateBodyElementCommon(data json.RawMessage, elementType string) []error {
	var ret []error
	var element bodyCommonTiniest
	err := json.Unmarshal(data, &element)
	if err != nil {
		ret = append(ret, fmt.Errorf("Unable to unmarshal element of type: %v", elementType))
	} else {
		// TODO: Maybe process paragraphs with empty content ?
		// if len(element.Content) == 0 {
		// 	ret = append(ret, fmt.Errorf("Mandatory field content is empty in element of type: %v", elementType))
		// }
	}
	return ret
}

func validateBodyElementImage(data json.RawMessage, elementType string) []error {
	var ret []error
	var element bodyImageTiniest
	err := json.Unmarshal(data, &element)
	if err != nil {
		ret = append(ret, fmt.Errorf("Unable to unmarshal element of type: %v", elementType))
	} else {
		if len(element.Source) == 0 {
			ret = append(ret, fmt.Errorf("Mandatory field source is empty in element of type: %v", elementType))
		}
	}
	return ret
}

func validateBodyElementVideo(data json.RawMessage, elementType string) []error {
	var ret []error
	var element bodyVideoTiniest
	err := json.Unmarshal(data, &element)
	if err != nil {
		ret = append(ret, fmt.Errorf("Unable to unmarshal element of type: %v", elementType))
	} else {
		if len(element.Source) == 0 {
			ret = append(ret, fmt.Errorf("Mandatory field source is empty in element of type: %v", elementType))
		}
		if len(element.VideoType) == 0 {
			ret = append(ret, fmt.Errorf("Mandatory field videoType is empty in element of type: %v", elementType))
		} else {
			if element.VideoType != supportedVideoTypeVimeo && element.VideoType != supportedVideoTypeYoutube {
				ret = append(ret, fmt.Errorf("Mandatory field videoType has invalid content (%v) in element of type: %v", element.VideoType, elementType))

			}
		}
	}
	return ret
}

func validateBodyElementGallery(data json.RawMessage, elementType string) []error {
	var ret []error
	var element bodyGalleryTiniest
	err := json.Unmarshal(data, &element)
	if err != nil {
		ret = append(ret, fmt.Errorf("Unable to unmarshal element of type: %v", elementType))
	} else {
		if len(element.Body) == 0 {
			ret = append(ret, fmt.Errorf("Mandatory field body is empty in element of type: %v", elementType))
		} else {
			ret = append(ret, validateBody(element.Body)...)
		}
	}
	return ret
}
