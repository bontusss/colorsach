package routes

import (
	"github.com/bontusss/colosach/controllers"
	"github.com/bontusss/colosach/middleware"
	"github.com/bontusss/colosach/services"
	"github.com/gin-gonic/gin"
)

type LibraryRouteController struct {
	libController controllers.LibraryController
}

func NewLibRouteController(libController controllers.LibraryController) LibraryRouteController {
	return LibraryRouteController{libController: libController}
}

func (l *LibraryRouteController) LibraryRoute(rg *gin.RouterGroup, userService services.UserService) {
	router := rg.Group("libs")
	router.Use(middleware.DeserializeUser(userService))
	router.POST("/", l.libController.CreateLibrary)
	router.GET("/", l.libController.FindLibraries)
}
