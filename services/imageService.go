package services

import (
	"mime/multipart"

	"github.com/bontusss/colosach/models"
)

type ImageService interface {
	UploadImage(*models.UploadImageInput, *multipart.FileHeader) (*models.ImageResponse, error)
}
