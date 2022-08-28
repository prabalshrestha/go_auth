package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// User structure
type User struct {
	ID           bson.ObjectId `json:"_id" bson:"_id"`
	Name         string        `json:"name" bson:"name"`
	Email        string        `json:"email" bson:"email"`
	Password     string        `json:"password" bson:"password" validate:"min=8,max=40,regexp=^[a-zA-Z]*$"`
	UserRole     []string      `json:"userRole" bson:"userRole"`
	UserId       string        `json:"userId" bson:"userId"`
	Token        string        `json:"userToken" bson:"userToken"`
	RefreshToken string        `json:"userRefreshToken" bson:"userRefreshToken"`
	CreatedAt    time.Time     `json:"createdAt" bson:"created_at"`
	UpdatedAt    time.Time     `json:"updatedAt" bson:"updated_at"`
}

// Users list
type Users []User
