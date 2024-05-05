package routes

import (
	"github.com/bontusss/colosach/controllers"
	"github.com/bontusss/colosach/middleware"
	"github.com/bontusss/colosach/services"
	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup, userService services.UserService) {

	router := rg.Group("user")
	router.Use(middleware.DeserializeUser(userService))
	// Get current user details
	router.GET("/me", uc.userController.GetMe)
	// make a user an admin
	router.PATCH("/make-admin", uc.userController.MakeAdminUser)
	// update a user
	router.POST("/update-me/:id", uc.userController.UpdateUser)
}
