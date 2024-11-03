package util

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetObjectById[T any](db *gorm.DB, id uuid.UUID) (*T, error) {
	obj := new(T)
	err := db.Model(obj).Where("id = ?", id.String()).First(&obj).Error

	return obj, err
}
