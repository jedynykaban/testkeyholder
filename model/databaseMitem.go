package model

import (
	"encoding/json"
	"time"
)

const (
	// StatusDelete indicates that a resource e.g. an article is meant for deletion.
	StatusDelete = 1
)

// TODO: Change the name after the refactoring
// DatabaseMitem represents the structure we have in a database
// that describes a mitem.
type DatabaseMitem struct {
	ID         string
	Data       json.RawMessage `datastore:",noindex"`
	SourceURL  string
	LogoURL    string
	SourceID   string
	Slug       string
	Status     int
	UserEdited bool
}

type MitemInPlaylist struct {
	ID           string
	CreationDate time.Time
	Inactive     bool
}
