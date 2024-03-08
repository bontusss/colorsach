package controllers

import (
	"github.com/bontusss/colosach/models"
	"github.com/bontusss/colosach/services"
	"github.com/bontusss/colosach/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type LibraryController struct {
	LibraryService services.LibraryService
}

func NewLibraryController(libService services.LibraryService) LibraryController {
	return LibraryController{libService}
}

func (l *LibraryController) CreateLibrary(c *gin.Context) {
	var lib *models.CreateLibraryRequest

	lib.User = utils.GetCurrentUser(c)
	if lib.User == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

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
