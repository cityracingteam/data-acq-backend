package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JwtKey struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid"`
	Data []byte
}
