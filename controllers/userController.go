package controllers

import (
	"log"
	"net/http"

	"github.com/bontusss/colosach/models"
	"github.com/bontusss/colosach/services"
	"github.com/bontusss/colosach/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return UserController{userService}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.DBResponse)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": models.FilteredResponse(currentUser)}})
}

// MakeAdminUser gets a takes a user ID and updates User.Role to UserRoleAdmin
// this can only be done by a user who has UserRoleSuperAdmin role
func (uc *UserController) MakeAdminUser(c *gin.Context) {
	userID := c.Params.ByName("user_id")
	currentUser := utils.GetCurrentUser(c)
	if currentUser.Role == models.UserRoleSuperAdmin {
		updateData := models.UpdateInput{Role: models.UserRoleAdmin}
		newAdminUser, err := uc.userService.UpdateUserById(userID, &updateData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "error creating admin user"})
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, newAdminUser)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Unauthorized action"})
	}
}
