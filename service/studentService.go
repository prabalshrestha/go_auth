package service

import (
	"login/conn"
	model "login/models"

	"gopkg.in/mgo.v2/bson"
)

type studentService interface {
	FindAllStudents(userCollection string) (model.Users, error)
	FindUserById(id bson.ObjectId, userCollection string) (model.User, error)
	AddUser(user model.User, userCollection string) error
	FindByEmail(email string, userCollection string) (model.User, error)
	UpdateUser(user model.User, userCollection string) error
	RemoveUser(id bson.ObjectId, userCollection string) error
	Exists(field string, value string, userCollection string) bool
}

type studentServiceImpl struct {
}

var (
	StudentService studentService = &studentServiceImpl{}
)

func (studentService *studentServiceImpl) FindAllStudents(userCollection string) (model.Users, error) {
	// Get DB from Mongo Config
	db := conn.GetMongoDB()
	users := model.Users{}
	err := db.C(userCollection).Find(bson.M{"userRole": "student"}).All(&users)
	return users, err
}
func (studentService *studentServiceImpl) FindUserById(id bson.ObjectId, userCollection string) (model.User, error) {
	db := conn.GetMongoDB()
	user := model.User{}
	err := db.C(userCollection).Find(bson.M{"_id": &id}).One(&user)
	return user, err
}

func (studentService *studentServiceImpl) AddUser(user model.User, userCollection string) error {
	db := conn.GetMongoDB()
	err := db.C(userCollection).Insert(user)
	return err
}

func (studentService *studentServiceImpl) FindByEmail(email string, userCollection string) (model.User, error) {
	db := conn.GetMongoDB()
	user := model.User{}
	err := db.C(userCollection).Find(bson.M{"email": email}).One(&user)
	return user, err
}
func (studentService *studentServiceImpl) UpdateUser(user model.User, userCollection string) error {
	db := conn.GetMongoDB()
	err := db.C(userCollection).Update(bson.M{"_id": user.ID}, user)
	return err
}

func (studentService *studentServiceImpl) RemoveUser(id bson.ObjectId, userCollection string) error {
	db := conn.GetMongoDB()
	err := db.C(userCollection).Remove(bson.M{"_id": &id})
	return err
}

func (studentService *studentServiceImpl) Exists(field string, value string, userCollection string) bool {
	db := conn.GetMongoDB()
	user := model.User{}
	err := db.C(userCollection).Find(bson.M{field: value}).One(&user)
	if err == nil {
		return true
	}
	return false
}
