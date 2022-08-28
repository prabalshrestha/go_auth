package controllers

import (
	"errors"
	"login/models"
	service "login/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

const RoleCollection = "role"

func CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.Bind(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errInvalidBody.Error()})
		return
	}
	if service.RoleService.Exists("name", role.Name, RoleCollection) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errors.New("role already exits")})
		return
	}
	role.ID = bson.NewObjectId()
	role.Selected = false
	err := service.RoleService.AddRole(role, RoleCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "role": &role})

}

func GetRole(c *gin.Context) {
	var id bson.ObjectId = bson.ObjectIdHex(c.Param("id"))
	role, err := service.RoleService.FindRoleById(id, RoleCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errInvalidID.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "role": &role})

}

func GetAllRoles(c *gin.Context) {
	// Get DB from Mongo Config
	roles, err := service.RoleService.FindAllRoles(RoleCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errNotExist.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "roles": &roles})
}

func UpdateRole(c *gin.Context) {
	var id bson.ObjectId = bson.ObjectIdHex(c.Param("id"))
	existingRole, err := service.RoleService.FindRoleById(id, RoleCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errInvalidID.Error()})
		return
	}
	err = c.Bind(&existingRole)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errInvalidBody.Error()})
		return
	}
	existingRole.ID = id
	err = service.RoleService.UpdateRole(existingRole, RoleCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errUpdationFailed.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "user": &existingRole})

}

func DeleteRole(c *gin.Context) {
	var id bson.ObjectId = bson.ObjectIdHex(c.Param("id")) // Get Param
	err := service.RoleService.RemoveRole(id, RoleCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "error": errDeletionFailed.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})

}
