package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	helper "login/helpers"
	model "login/models"
	service "login/service"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

type Response struct {
	Status string     `json:"status"`
	User   model.User `json:"user"`
	Error  string     `json:"error"`
}

var (
	mockUser model.User = model.User{
		ID:        "61af691593fa4a2a7bdc43ec",
		Name:      "testUser",
		Email:     "test@gmail.com",
		Password:  "$2a$14$.39kjQYyB3T6x.0b8NDRG.6sqGpUMfsTKHe/uDyQa98p7hMGBAeg.",
		UserId:    "61af691593fa4a2a7bdc43ec",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	mockUsers model.Users = model.Users{
		{
			ID:        "61af691593fa4a2a7bdc43ec",
			Name:      "testname",
			Email:     "u@gmail.com",
			Password:  "$2a$14$.39kjQYyB3T6x.0b8NDRG.6sqGpUMfsTKHe/uDyQa98p7hMGBAeg.",
			UserId:    "61af691593fa4a2a7bdc43ec",
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
	}
)

type userServiceMock struct {
	handleAddUserFn    func() error
	handleFindAllUsers func() (model.Users, error)
	handleFindUserById func() (model.User, error)
	handleFindByEmail  func() (model.User, error)
	handleUpdateUser   func() error
	handleRemoveUser   func() error
	handleExists       func() bool
}

func (mock *userServiceMock) FindAllUsers(userCollection string) (model.Users, error) {
	return mock.handleFindAllUsers()
}
func (mock *userServiceMock) FindUserById(id bson.ObjectId, userCollection string) (model.User, error) {
	return mock.handleFindUserById()
}
func (mock *userServiceMock) AddUser(user model.User, userCollection string) error {
	return mock.handleAddUserFn()
}
func (mock *userServiceMock) FindByEmail(email string, userCollection string) (model.User, error) {
	return mock.handleFindByEmail()
}
func (mock *userServiceMock) UpdateUser(user model.User, userCollection string) error {
	return mock.handleUpdateUser()
}
func (mock *userServiceMock) RemoveUser(id bson.ObjectId, userCollection string) error {
	return mock.handleRemoveUser()
}
func (mock *userServiceMock) Exists(field string, value string, userCollection string) bool {
	return mock.handleExists()
}

// MockPasswordhelper
type passwordHelperMock struct {
	handleHashPassword           func() (string, error)
	handleVerifyPassword         func() (bool, string)
	handleGenerateRandomPassword func() string
}

func (mock *passwordHelperMock) HashPassword(password string) (string, error) {
	return mock.handleHashPassword()
}
func (mock *passwordHelperMock) VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	return mock.handleVerifyPassword()
}
func (mock *passwordHelperMock) GenerateRandomPassword() string {
	return mock.handleGenerateRandomPassword()
}

// LoginTest{

func TestLoginSuccess(t *testing.T) {
	serviceMock := userServiceMock{}
	serviceMock.handleFindByEmail = func() (model.User, error) {
		return mockUser, nil
	}
	serviceMock.handleUpdateUser = func() error {
		return nil
	}
	service.UserService = &serviceMock

	passwdHeplerMock := passwordHelperMock{}
	passwdHeplerMock.handleVerifyPassword = func() (bool, string) {
		check := true
		var msg string = ""
		return check, msg
	}
	helper.PasswordHelper = &passwdHeplerMock
	var jsonStr = []byte(`
		{
		"email":"test@gmail.com",
		"passsword":"test"
		}
		`)
	response := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/api/users/login", bytes.NewBuffer(jsonStr))

	c, _ := gin.CreateTestContext(response)
	c.Request = request
	fmt.Println(c.Request)
	Login(c)
	var jsonRes Response
	_ = json.Unmarshal(response.Body.Bytes(), &jsonRes)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "success", jsonRes.Status)
}

func TestLoginInvalidBody(t *testing.T) {
	serviceMock := userServiceMock{}
	serviceMock.handleFindByEmail = func() (model.User, error) {
		return mockUser, nil
	}
	serviceMock.handleUpdateUser = func() error {
		return nil
	}
	service.UserService = &serviceMock

	passwdHeplerMock := passwordHelperMock{}
	passwdHeplerMock.handleVerifyPassword = func() (bool, string) {
		check := true
		var msg string = ""
		return check, msg
	}
	helper.PasswordHelper = &passwdHeplerMock
	var jsonStr = []byte(`{}`)
	response := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/api/users/login", bytes.NewBuffer(jsonStr))

	c, _ := gin.CreateTestContext(response)
	c.Request = request
	Login(c)
	var jsonRes Response
	_ = json.Unmarshal(response.Body.Bytes(), &jsonRes)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, "failed", jsonRes.Status)
	assert.Equal(t, "invalid request body", jsonRes.Error)
}
func TestLoginInvalidEmail(t *testing.T) {
	serviceMock := userServiceMock{}
	serviceMock.handleFindByEmail = func() (model.User, error) {
		return mockUser, errInvalidCred
	}
	serviceMock.handleUpdateUser = func() error {
		return nil
	}
	service.UserService = &serviceMock

	passwdHeplerMock := passwordHelperMock{}
	passwdHeplerMock.handleVerifyPassword = func() (bool, string) {
		check := true
		var msg string = ""
		return check, msg
	}
	helper.PasswordHelper = &passwdHeplerMock
	var jsonStr = []byte(`{
		"email":"test@gmail.com",
		"passsword":"test"}`)
	response := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/users/signup", bytes.NewBuffer(jsonStr))

	c, _ := gin.CreateTestContext(response)
	c.Request = request
	Login(c)
	var jsonRes Response
	_ = json.Unmarshal(response.Body.Bytes(), &jsonRes)
	fmt.Println(jsonRes.Error, jsonRes.Status)
	assert.Equal(t, 401, response.Code)
	assert.Equal(t, "failed", jsonRes.Status)
	assert.Equal(t, "invalid email or password", jsonRes.Error)
}
func TestLoginInvalidPassword(t *testing.T) {
	serviceMock := userServiceMock{}
	serviceMock.handleFindByEmail = func() (model.User, error) {
		return mockUser, nil
	}
	serviceMock.handleUpdateUser = func() error {
		return nil
	}
	service.UserService = &serviceMock

	passwdHeplerMock := passwordHelperMock{}
	passwdHeplerMock.handleVerifyPassword = func() (bool, string) {
		check := false
		var msg string = ""
		return check, msg
	}
	helper.PasswordHelper = &passwdHeplerMock
	var jsonStr = []byte(`{
		"email":"test@gmail.com",
		"passsword":"test"}`)
	response := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "api/users/login", bytes.NewBuffer(jsonStr))

	c, _ := gin.CreateTestContext(response)
	c.Request = request
	Login(c)
	var jsonRes Response
	_ = json.Unmarshal(response.Body.Bytes(), &jsonRes)
	fmt.Println(jsonRes.Error, jsonRes.Status)
	assert.Equal(t, 401, response.Code)
	assert.Equal(t, "failed", jsonRes.Status)
	assert.Equal(t, "invalid email or password", jsonRes.Error)
}
func TestLoginUpdateFailed(t *testing.T) {
	serviceMock := userServiceMock{}
	serviceMock.handleFindByEmail = func() (model.User, error) {
		return mockUser, nil
	}
	serviceMock.handleUpdateUser = func() error {
		return errUpdationFailed
	}
	service.UserService = &serviceMock

	passwdHeplerMock := passwordHelperMock{}
	passwdHeplerMock.handleVerifyPassword = func() (bool, string) {
		check := true
		var msg string = ""
		return check, msg
	}
	helper.PasswordHelper = &passwdHeplerMock
	var jsonStr = []byte(`{
		"email":"test@gmail.com",
		"passsword":"test"}`)
	response := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/users/signup", bytes.NewBuffer(jsonStr))

	c, _ := gin.CreateTestContext(response)
	c.Request = request
	Login(c)
	var jsonRes Response
	_ = json.Unmarshal(response.Body.Bytes(), &jsonRes)
	fmt.Println(jsonRes.Error, jsonRes.Status)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, "failed", jsonRes.Status)
	assert.Equal(t, "error in the user updation", jsonRes.Error)
}

func TestCreateUserSuccess(t *testing.T) {
	serviceMock := userServiceMock{}
	serviceMock.handleAddUserFn = func() error {
		return nil
	}
	service.UserService = &serviceMock

	var jsonStr = []byte(`{
				"name":"testUser",
				"email":"test@gmail.com",
				"address":"test",
				"passsword":"test",
				"random":"random",
				"age":21}`)
	response := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/users/signup", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(response)
	c.Request = request
	CreateUser(c)
	var jsonRes Response
	_ = json.Unmarshal(response.Body.Bytes(), &jsonRes)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "success", jsonRes.Status)
}

func TestCreateUsererrInsertionFailed(t *testing.T) {
	mockErr := errors.New("mocked Err")
	serviceMock := userServiceMock{}
	serviceMock.handleAddUserFn = func() error {
		return mockErr
	}
	service.UserService = &serviceMock

	var jsonStr = []byte(`{
				"name":"testUser",
				"email":"test@gmail.com",
				"address":"test",
				"passsword":"test",
				"random":"random",
				"age":21}`)
	response := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/users/signup", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(response)
	c.Request = request
	CreateUser(c)
	var jsonRes Response
	_ = json.Unmarshal(response.Body.Bytes(), &jsonRes)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, "failed", jsonRes.Status)
}
func TestCreateUserInvalidBody(t *testing.T) {
	serviceMock := userServiceMock{}
	serviceMock.handleAddUserFn = func() error {
		return nil
	}
	service.UserService = &serviceMock

	var jsonStr = []byte(`{
				"name":"testUser",
				"email":"test@gmail.com",
				"address":"test",
				"passsword":"test",
				"random":"random"
				"age":21}`)
	response := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/users/signup", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(response)
	c.Request = request
	CreateUser(c)
	var jsonRes Response
	_ = json.Unmarshal(response.Body.Bytes(), &jsonRes)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, "failed", jsonRes.Status)
	assert.Equal(t, "invalid request body", jsonRes.Error)
}
