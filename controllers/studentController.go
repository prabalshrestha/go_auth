package controllers

import (
	"log"
	helper "login/helpers"
	model "login/models"
	service "login/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// GetAllUser Endpoint
func GetAllStudents(c *gin.Context) {
	// Get DB from Mongo Config
	users, err := service.StudentService.FindAllStudents(UserCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errNotExist.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "users": &users})
}

// SIGNUP
func RegisterStudent(c *gin.Context) {
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
	password, _ := helper.PasswordHelper.HashPassword(user.Password)
	user.Password = password
	user.UserRole = append(user.UserRole, "student")
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
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   1000,
		HttpOnly: true,
		Path:     "/",
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		MaxAge:   1000,
		HttpOnly: true,
		Path:     "/",
	})
	c.Cookie("token")
	err = service.UserService.AddUser(user, UserCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errInsertionFailed.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "user": &user})
}
