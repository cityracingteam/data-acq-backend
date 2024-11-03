package models

import (
	"github.com/google/uuid"
	"github.com/markbates/goth"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid"`
	FirstName string
	LastName  string
	Email     string
	AvatarURL string
}

func GetUserFromGoth(g *goth.User) *User {
	uuid, err := uuid.Parse(g.UserID)
	if err != nil {
		// return an empty user object on parse error...
		// TODO: return error object
		return &User{}
	}
	user := &User{
		ID:        uuid,
		FirstName: g.FirstName,
		LastName:  g.LastName,
		Email:     g.Email,
		AvatarURL: g.AvatarURL,
	}
	return user
}
