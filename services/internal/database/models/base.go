package models

import (
	"time"

	"github.com/google/uuid"
)

// Base represents the base class
type Base struct {
	uid          uuid.UUID
	dateCreated  time.Time
	dateModified time.Time
}

// NewBase creates a new instance of Base
func NewBase() *Base {
	return &Base{
		uid:          uuid.New(),
		dateCreated:  time.Now(),
		dateModified: time.Now(),
	}
}

// ID returns the ID of the base object
func (b *Base) ID() uuid.UUID {
	return b.uid
}

// DateCreated returns the creation date of the base object
func (b *Base) DateCreated() time.Time {
	return b.dateCreated
}

// DateModified returns the modification date of the base object
func (b *Base) DateModified() time.Time {
	return b.dateModified
}
