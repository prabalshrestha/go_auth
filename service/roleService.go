package service

import (
	"login/conn"
	model "login/models"

	"gopkg.in/mgo.v2/bson"
)

type roleService interface {
	FindAllRoles(Collection string) (model.Roles, error)
	FindRoleById(id bson.ObjectId, Collection string) (model.Role, error)
	AddRole(role model.Role, Collection string) error
	UpdateRole(role model.Role, Collection string) error
	RemoveRole(id bson.ObjectId, Collection string) error
	Exists(field string, value string, Collection string) bool
}

type roleServiceImpl struct {
}

var (
	RoleService roleService = &roleServiceImpl{}
)

func (roleService *roleServiceImpl) FindAllRoles(Collection string) (model.Roles, error) {
	// Get DB from Mongo Config
	db := conn.GetMongoDB()
	roles := model.Roles{}
	err := db.C(Collection).Find(bson.M{}).All(&roles)
	return roles, err
}
func (roleService *roleServiceImpl) FindRoleById(id bson.ObjectId, Collection string) (model.Role, error) {
	db := conn.GetMongoDB()
	role := model.Role{}
	err := db.C(Collection).Find(bson.M{"_id": &id}).One(&role)
	return role, err
}

func (roleService *roleServiceImpl) AddRole(role model.Role, Collection string) error {
	db := conn.GetMongoDB()
	err := db.C(Collection).Insert(role)
	return err
}

func (roleService *roleServiceImpl) UpdateRole(role model.Role, Collection string) error {
	db := conn.GetMongoDB()
	err := db.C(Collection).Update(bson.M{"_id": role.ID}, role)
	return err
}

func (roleService *roleServiceImpl) RemoveRole(id bson.ObjectId, Collection string) error {
	db := conn.GetMongoDB()
	err := db.C(Collection).Remove(bson.M{"_id": &id})
	return err
}

func (roleService *roleServiceImpl) Exists(field string, value string, Collection string) bool {
	db := conn.GetMongoDB()
	role := model.Role{}
	err := db.C(Collection).Find(bson.M{field: value}).One(&role)
	if err == nil {
		return true
	}
	return false
}
