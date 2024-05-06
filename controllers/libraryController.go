package controllers

import (
	"fmt"
	"net/http"

	"github.com/bontusss/colosach/models"
	"github.com/bontusss/colosach/services"
	"github.com/bontusss/colosach/utils"

	"github.com/gin-gonic/gin"
)

type LibraryController struct {
	LibraryService services.LibraryService
}

func NewLibraryController(libService services.LibraryService) LibraryController {
	return LibraryController{libService}
}

func (l *LibraryController) CreateLibrary(c *gin.Context) {
	lib := &models.CreateLibraryRequest{}
	// check if user is logged in
	currentUser := utils.GetCurrentUser(c)
	lib.OwnerID = currentUser.ID.Hex()

	if err := c.ShouldBindJSON(&lib); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "an error occurred, try again."})
		fmt.Println("creating library error: ", err)
		return
	}

	newLib, err := l.LibraryService.CreateLibrary(lib)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": "error creating library."})
		fmt.Println("error creating library: ", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": newLib})
}
