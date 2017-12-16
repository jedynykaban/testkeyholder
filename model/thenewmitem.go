// TODO: change the package name to mitem
// once we will get rid of the Mitem structure in jukebox
package model

import (
	"strconv"
	"strings"
	"time"

	"encoding/json"

	log "github.com/Sirupsen/logrus"
)

// TODO: change the name to Mitem
// once we will get rid of the Mitem structure in jukebox
// TODO: find a better place for the package
type TheNewMitem struct {
	Type         string    `json:"type"`
	ID           string    `json:"id"`
	Headline     string    `json:"headline"`
	Slug         string    `json:"slug"`
	MainImage    Image     `json:"mainImage"`
	CreationDate time.Time `json:"creationDate"`
	Status       int       `json:"status"`
	Body         Body      `json:"body"`
	Meta         Meta      `json:"meta,omitempty"`
}

// Image defines an image structure
type Image struct {
	Source  string `json:"source"`
	Caption string `json:"caption"`
	Height  int    `json:"height"`
	Width   int    `json:"width"`
}

// Meta represents meta data attached to a mitem
type Meta struct {
	// SourceID the ID of the initial playlist
	// the mitem was added to.
	SourceID string `json:"sourceID"`

	// SourceURL holds origin URL i.e. https://bbc.co.uk
	SourceURL string `json:"sourceURL"`

	// At the moment used to display the logo of the mgazine in the rendered article
	// TODO: LogoURL is magazine specific thus should not be part of Meta
	LogoURL string `json:"logoURL"`

	MosaiqPrimary MosaiqPrimary `json:"mosaiqPrimary"`

	// UserEdited indicates whether a mitem was edited manually
	// by a user (via the admin tool)
	UserEdited bool `json:"userEdited"`

	// Inactive indicates whether a mitem should be included in a playlist
	// By default playlists don't contain inactive mitems
	Inactive bool `json:"inactive"`

	// TODO: change to strongly typed structure along with some payload.
	License interface{} `json:"license"`

	// Analytics a hint to the client to send metricts to another destination
	Analytics []Analytics `json:"analytics"`

	// Authors contains the list of authors along with theirs playlist's name
	Authors []Author `json:"authors"`

	// SectionPlaylist tells the client that we have more than one section in a magazine
	// TODO: SectionPlaylist is magazine specific thus should not be part of Meta
	SectionPlaylist []SectionPlaylist `json:"sectionPlaylist"`

	// Section formerly known as category, holds the information about
	// taxonomy. It also assings a colour to each group.
	// TODO: Section is magazine specific thus should not be part of Meta
	Section Section `json:"section"`

	// AdsPolicy describes ads policy that should be enforced by the client
	// i.e. maximum number of ads in an article.
	AdsPolicy AdsPolicy `json:"adsPolicy"`

	// Tags represents a collection of tags that were attached to a mitem
	Tags []Tag `json:"tags, omitempty"`
}

// Tag represents a tag attached to a mitem
type Tag struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// Analytics a hint to the client to send metricts to another destination
type Analytics struct {
	// Type holds type information i.e. ga (Google Analytics)
	Type string `json:"type"`
	// ID e.g. UA-89047834-1 (GA tracing ID)
	ID string `json:"ID"`
}

// MosaiqPrimary means that an article was
// created inside mosaiq. As opposed to pulling from
// an external sources.
// On the client we set <link rel="canonical" href="http://skonahem.com/?p=64187"> to some value
type MosaiqPrimary struct {
	// Set indicates if "MosaiqPrimary" was set
	Set    bool   `json:"set"`
	Domain string `json:"domain"`
}

// Section holds the information about colours and taxonomy
type Section struct {
	Colour   Colour `json:"color"`
	Gradient Colour `json:"gradient"`
	Tier1    string `json:"tier1"`
	Tier2    string `json:"tier2"`
}

// Colour represents a colour
type Colour struct {
	Hue        float64 `json:"h"`
	Saturation float64 `json:"s"`
	Lum        float64 `json:"l"`
}

var blankColour = Colour{0, 0, 0}

// NewColour creates new Color from a string.
// It expects that th input string will be separeted by |.
func NewColour(c string) Colour {
	if len(c) == 0 {
		return blankColour
	}

	minS := 3
	cs := strings.Split(c, "|")
	if len(cs) < minS {
		return blankColour
	}

	hue, err := strconv.ParseFloat(cs[0], 64)
	if err != nil {
		log.Errorf("cannot convert color hue into numeric")
		return blankColour
	}

	sat, err := strconv.ParseFloat(cs[1], 64)
	if err != nil {
		log.Errorf("cannot convert color sat into numeric")
		return blankColour
	}

	lum, err := strconv.ParseFloat(cs[2], 64)
	if err != nil {
		log.Errorf("cannot convert color lum into numeric")
		return blankColour
	}
	return Colour{hue, sat, lum}
}

// AdsPolicy describes ads policy that should be enforced by the client
// i.e. maximum number of ads in an article.
type AdsPolicy struct {
	On     bool `json:"on"`
	MaxAds int  `json:"maxAds"`
}

type SectionPlaylist struct {
	DisplayName  string `json:"displayName"`
	PlaylistName string `json:"playlistName"`
}

// Author holds the information about an article's author
type Author struct {
	Name         string `json:"name"`
	PlaylistName string `json:"playlistName"`
}

// Body contains serialized body elements like paragraphs, tables, lists, images, etc.
type Body []json.RawMessage
