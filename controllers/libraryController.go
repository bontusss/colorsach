package controllers

import (
	"fmt"
	"net/http"
	"strconv"

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


// @Summary Create library
// @Description Create a library
// @Tags Library
// @Accept json
// @Produce json
// @Param DBLibrary body models.DBLibrary true "DBLibrary"
// @Success 201
// @Failure 500
// @Failure 400
// @Router /api/libs [post]
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
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		fmt.Println("error creating library: ", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": newLib})
}

func (l *LibraryController) FindLibraries(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	libraries, err := l.LibraryService.FindLibraries(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": libraries})
}