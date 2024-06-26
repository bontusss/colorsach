package controllers

import (
	"log"
	"net/http"

	"github.com/bontusss/colosach/models"
	"github.com/bontusss/colosach/services"
	"github.com/bontusss/colosach/utils"
	"github.com/gin-gonic/gin"
)

type ImageController struct {
	imageService services.ImageService
}

func NewImageController(imageService services.ImageService) ImageController {
	return ImageController{imageService: imageService}
}

func (i *ImageController) UploadImage(c *gin.Context) {
	user := utils.GetCurrentUser(c)
	imagee := &models.UploadImageInput{}
	imagee.Title = c.PostForm("title")
	imagee.Tags = c.PostFormArray("tags")

	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "File too large"})
		return
	}
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		log.Println("An error occurred", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve file"})
		return
	}
	defer file.Close()

	if header.Size == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Empty file"})
		return
	}

	imageData, err := i.imageService.UploadImage(imagee, file, user)
	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "success", "data": imageData})
}
