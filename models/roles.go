package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Role structure
type Role struct {
	ID          bson.ObjectId `json:"_id" bson:"_id"`
	Name        string        `json:"name" bson:"name"`
	Selected    bool          `json:"selected" bson:"selected"`
	Permissions []string      `json:"permission" bson:"permission"`
}

// Roles list
type Roles []Role
