package controllers

import (
	"log"
	"net/http"
	"time"

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

// MakeAdminUser gets a takes a user ID and updates User.Role to UserRoleAdmin
// this can only be done by a user who has UserRoleSuperAdmin role
// @Summary Get Current User
// @Description Get the details of a logged in user
// @Tags User
// @Accept json
// @Produce json
// @Param DBResponse body models.DBResponse true "DBResponse"
// @Success 200
// @Failure 401
// @Router /api/users/me [get]
func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.DBResponse)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": models.FilteredResponse(currentUser)}})
}



// @Summary Update user profile
// @Description Users update their profile
// @Tags User
// @Accept json
// @Produce json
// @Param UserResponse body models.UserResponse true "UserResponse"
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 500
// @Router /api/users/update-me/:id [post]
func (uc *UserController) UpdateUser(c *gin.Context) {
	var data *models.UpdateInput
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "invalid request body"})
		return
	}

	currentUser := utils.GetCurrentUser(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "user not authenticated"})
		return
	}
	// Ensure users can only update their own details
	if c.Param("id") != currentUser.ID.Hex() {
		c.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "cannot update other user's details"})
		return
	}
	updatedUser, err := uc.userService.UpdateUserById(currentUser.ID.Hex(), data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedUser})
}

type data struct {
	Email string `json:"email"`
}

// MakeAdminUser gets a takes a user ID and updates User.Role to UserRoleAdmin
// this can only be done by a user who has UserRoleSuperAdmin role
// @Summary Make Admin
// @Description Make a user an admin
// @Tags User
// @Accept json
// @Produce json
// @Param data body data true "data"
// @Success 200
// @Failure 400
// @Failure 401
// @Router /api/users/make-admin [patch]
func (uc *UserController) MakeAdminUser(c *gin.Context) {
	var reqData data
	if err := c.BindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "invalid request body"})
		return
	}
	currentUser := utils.GetCurrentUser(c)
	if currentUser.Role == models.UserRoleSuperAdmin {
		updateData := &models.UpdateInput{Role: models.UserRoleAdmin, UpdatedAt: time.Now()}
		// newAdminUser := &models.DBResponse{}
		newAdminUser, err := uc.userService.UpdateUserByEmail(reqData.Email, updateData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "error creating admin user"})
			log.Fatal("error creating admin user: ", err)
		}
		c.JSON(http.StatusOK, newAdminUser)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Unauthorized action"})
	}
}
