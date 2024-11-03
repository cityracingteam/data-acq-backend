package resolver

import (
	"log"

	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB *gorm.DB
}

// Set DB object if not already set
func (r *Resolver) UpdateDb(db *gorm.DB) {
	if r.DB == nil {
		log.Printf("[resolver/UpdateDb]: db not set, setting.")
		r.DB = db
	} else {
		log.Printf("[resolver/UpdateDb]: WARN: db is already set, not overwriting.")
	}
}
