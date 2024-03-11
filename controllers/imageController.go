package controllers

import (
	"log"
	"net/http"

	"github.com/bontusss/colosach/models"
	"github.com/bontusss/colosach/services"
	"github.com/gin-gonic/gin"
)

type ImageController struct {
	imageService services.ImageService
}

func NewImageController(imageService services.ImageService) ImageController {
	return ImageController{imageService: imageService}
}

func (i *ImageController) UploadImage(c *gin.Context) {
	// Parse from form data
	image := &models.UploadImageInput{}
	image.Name = c.PostForm("name")
	image.Tags = c.PostFormArray("tags")
	file, err := c.FormFile("image")

	if err != nil {
		log.Println("An error occurred", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
	}

	imageData, err := i.imageService.UploadImage(image, file)
	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	// image.Image_url = url
	c.JSON(200, gin.H{"status": "success", "data": imageData})
}
