package controllers

import (
	"net/http"
	"strings"

	"github.com/bontusss/colosach/models"
	"github.com/bontusss/colosach/services"

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

	currentUser := c.MustGet("currentUser").(*models.DBResponse)
	libOwner := models.FilteredResponse(currentUser)
	lib.Owner = &libOwner

	if err := c.ShouldBindJSON(&lib); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newLib, err := l.LibraryService.CreateLibrary(lib)
	if err != nil {
		if strings.Contains(err.Error(), "title already exists") {
			c.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": newLib})
}
