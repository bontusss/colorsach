package routes

import (
	"github.com/bontusss/colosach/controllers"
	"github.com/bontusss/colosach/middleware"
	"github.com/bontusss/colosach/services"
	"github.com/gin-gonic/gin"
)

type ImageRouteController struct {
	imageController controllers.ImageController
}

func NewImageRouteController(imageController controllers.ImageController) ImageRouteController {
	return ImageRouteController{imageController}
}

func (i *ImageRouteController) ImageRoute(rg *gin.RouterGroup, userService services.UserService) {
	router := rg.Group("image")
	router.Use(middleware.DeserializeUser(userService))
	router.POST("/upload", i.imageController.UploadImage)
}