package controllers

import (
	"errors"
	"log"
	helper "login/helpers"
	"login/middlewares"
	model "login/models"
	service "login/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"gopkg.in/mgo.v2/bson"
)

// UserCollection statically declared
const UserCollection = "user"

var (
	errNotExist        = errors.New("users doesnt not exist")
	errInvalidID       = errors.New("invalid id")
	errInvalidBody     = errors.New("invalid request body")
	errBlankFields     = errors.New("email and password cannot be empty")
	errInvalidCred     = errors.New("invalid email or password")
	errInvalidEmail    = errors.New("account doesnt exists")
	errInsertionFailed = errors.New("error in the insertion")
	errUpdationFailed  = errors.New("error in the updation")
	errDeletionFailed  = errors.New("error in the deletion")
	errSmtp            = errors.New("user created but error in smtp")
)

type changePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

// Login
func Login(c *gin.Context) {

	var user model.User
	var foundUser model.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errInvalidBody.Error()})
		return
	}

	if user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errBlankFields.Error()})
		return
	}
	foundUser, err := service.UserService.FindByEmail(user.Email, UserCollection)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "error": errInvalidCred.Error()})
		return
	}
	passwordIsValid, _ := helper.PasswordHelper.VerifyPassword(user.Password, foundUser.Password)
	if !passwordIsValid {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "error": errInvalidCred.Error()})
		return
	}
	token, refreshToken, _ := helper.TokenHelper.GenerateAllToken(foundUser.Name, foundUser.UserId, foundUser.UserRole)
	foundUser.UpdatedAt = time.Now()
	foundUser.Token = token
	foundUser.RefreshToken = refreshToken
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   20,
		HttpOnly: true,
		Path:     "/",
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		MaxAge:   20,
		HttpOnly: true,
		Path:     "/",
	})
	c.Cookie("token")
	err = service.UserService.UpdateUser(foundUser, UserCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errUpdationFailed.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "user": &foundUser})

}

// SIGNUP
func CreateUser(c *gin.Context) {
	user := model.User{}
	err := c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errInvalidBody.Error()})
		return
	}
	isValid, err := helper.ValidationHepler.ValidateNewUser(user)
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err.Error()})
		return
	}
	enteredPassword := user.Password
	password, _ := helper.PasswordHelper.HashPassword(user.Password)
	user.Password = password
	user.ID = bson.NewObjectId()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.UserId = bson.ObjectId(user.ID).Hex()
	if err != nil {
		log.Panic(err)
	}
	token, refreshToken, _ := helper.TokenHelper.GenerateAllToken(user.Name, user.UserId, user.UserRole)
	user.Token = token
	user.RefreshToken = refreshToken
	err = service.UserService.AddUser(user, UserCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errInsertionFailed.Error()})
		return
	}
	err = helper.MailHelper.SendRegistrationMail(user, enteredPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errSmtp.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "user": &user})
}

//LOGOUT
func Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
	})
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetActiveUser
func GetActiveUser(c *gin.Context) {
	var id bson.ObjectId = bson.ObjectIdHex(c.GetString("userId"))
	activeUser, err := service.UserService.FindUserById(id, UserCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errNotExist.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "user": &activeUser})
}

// GetAllUser Endpoint
func GetAllUser(c *gin.Context) {
	// Get DB from Mongo Config
	users, err := service.UserService.FindAllUsers(UserCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errNotExist.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "users": &users})
}

// GetUser Endpoint
func GetUser(c *gin.Context) {
	var id bson.ObjectId = bson.ObjectIdHex(c.Param("id")) // Get Param
	user, err := service.UserService.FindUserById(id, UserCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errInvalidID.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "user": &user})
}

// PUT
// UpdateUser Endpoint
func UpdateUser(c *gin.Context) {

	var id bson.ObjectId = bson.ObjectIdHex(c.Param("id")) // Get Param

	if !(middlewares.Contains(c.GetStringSlice("userRole"), "admin") || (c.GetString("userId") == c.Param("id"))) {
		c.JSON(http.StatusForbidden, gin.H{"status": "failed", "error": errors.New("Forbidden").Error()})
	}
	existingUser, err := service.UserService.FindUserById(id, UserCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errInvalidID.Error()})
		return
	}
	// user := model.User{}
	err = c.Bind(&existingUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errInvalidBody.Error()})
		return
	}
	existingUser.ID = id
	existingUser.UpdatedAt = time.Now()
	err = service.UserService.UpdateUser(existingUser, UserCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errUpdationFailed.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "user": &existingUser})
}

func ChangePassword(c *gin.Context) {
	var id bson.ObjectId = bson.ObjectIdHex(c.Param("id")) // Get Param
	if !(c.GetString("userId") == c.Param("id")) {
		c.JSON(http.StatusForbidden, gin.H{"status": "failed", "error": errors.New("Forbidden").Error()})
	}
	changeRequest := changePasswordRequest{}
	err := c.Bind(&changeRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errInvalidBody.Error()})
		return
	}
	existingUser, err := service.UserService.FindUserById(id, UserCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errInvalidID.Error()})
		return
	}
	passwordIsValid, _ := helper.PasswordHelper.VerifyPassword(changeRequest.CurrentPassword, existingUser.Password)
	if !passwordIsValid {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "error": errors.New("incorrect password").Error()})
		return
	}
	hashedPassword, _ := helper.PasswordHelper.HashPassword(changeRequest.NewPassword)
	existingUser.Password = hashedPassword
	existingUser.UpdatedAt = time.Now()
	err = service.UserService.UpdateUser(existingUser, UserCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errUpdationFailed.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "user": &existingUser})

}

func ResetPassword(c *gin.Context) {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errInvalidBody.Error()})
		return
	}

	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errBlankFields.Error()})
		return
	}
	foundUser, err := service.UserService.FindByEmail(user.Email, UserCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errInvalidEmail.Error()})
		return
	}
	oldPassword := foundUser.Password
	newPassword := helper.PasswordHelper.GenerateRandomPassword()
	password, _ := helper.PasswordHelper.HashPassword(newPassword)
	foundUser.Password = password
	foundUser.UpdatedAt = time.Now()
	err = service.UserService.UpdateUser(foundUser, UserCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errUpdationFailed.Error()})
		return
	}
	err = helper.MailHelper.SendResetPasswordMail(foundUser, newPassword)
	if err != nil {
		foundUser.Password = oldPassword
		_ = service.UserService.UpdateUser(foundUser, UserCollection)
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errSmtp.Error()})

		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "user": &foundUser})

}

// DELETE
// DeleteUser Endpoint
func DeleteUser(c *gin.Context) {
	var id bson.ObjectId = bson.ObjectIdHex(c.Param("id")) // Get Param
	err := service.UserService.RemoveUser(id, UserCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errDeletionFailed.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
