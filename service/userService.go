package service

import (
	"login/conn"
	model "login/models"

	"gopkg.in/mgo.v2/bson"
)

type userService interface {
	FindAllUsers(userCollection string) (model.Users, error)
	FindUserById(id bson.ObjectId, userCollection string) (model.User, error)
	AddUser(user model.User, userCollection string) error
	FindByEmail(email string, userCollection string) (model.User, error)
	UpdateUser(user model.User, userCollection string) error
	RemoveUser(id bson.ObjectId, userCollection string) error
	Exists(field string, value string, userCollection string) bool
}

type userServiceImpl struct {
}

var (
	UserService userService = &userServiceImpl{}
)

func (userService *userServiceImpl) FindAllUsers(userCollection string) (model.Users, error) {
	// Get DB from Mongo Config
	db := conn.GetMongoDB()
	users := model.Users{}
	err := db.C(userCollection).Find(bson.M{}).All(&users)
	return users, err
}
func (userService *userServiceImpl) FindUserById(id bson.ObjectId, userCollection string) (model.User, error) {
	db := conn.GetMongoDB()
	user := model.User{}
	err := db.C(userCollection).Find(bson.M{"_id": &id}).One(&user)
	return user, err
}

func (userService *userServiceImpl) AddUser(user model.User, userCollection string) error {
	db := conn.GetMongoDB()
	err := db.C(userCollection).Insert(user)
	return err
}

func (userService *userServiceImpl) FindByEmail(email string, userCollection string) (model.User, error) {
	db := conn.GetMongoDB()
	user := model.User{}
	err := db.C(userCollection).Find(bson.M{"email": email}).One(&user)
	return user, err
}
func (userService *userServiceImpl) UpdateUser(user model.User, userCollection string) error {
	db := conn.GetMongoDB()
	err := db.C(userCollection).Update(bson.M{"_id": user.ID}, user)
	return err
}

func (userService *userServiceImpl) RemoveUser(id bson.ObjectId, userCollection string) error {
	db := conn.GetMongoDB()
	err := db.C(userCollection).Remove(bson.M{"_id": &id})
	return err
}

func (userService *userServiceImpl) Exists(field string, value string, userCollection string) bool {
	db := conn.GetMongoDB()
	user := model.User{}
	err := db.C(userCollection).Find(bson.M{field: value}).One(&user)
	if err == nil {
		return true
	}
	return false
}
