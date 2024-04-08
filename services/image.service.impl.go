package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/Edlinorg/prominentcolor"
	"github.com/bontusss/colosach/models"
	"github.com/bontusss/colosach/utils"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ImageServiceImpl struct {
	ImageCollection *mongo.Collection
	ctx             context.Context
}

func NewImageService(imageCollection *mongo.Collection, ctx context.Context) ImageService {
	return &ImageServiceImpl{imageCollection, ctx}
}

// UploadImage implements ImageService.
func (i *ImageServiceImpl) UploadImage(img *models.UploadImageInput, file *multipart.FileHeader) (*models.ImageResponse, error) {
	// check if the name contains an color
	colors := []string{"black", "red", "blue"}
	for _, color := range colors {
		if strings.Contains(strings.ToLower(img.Name), color) {
			// add the color to tags
			img.Tags = append(img.Tags, color)
		}
	}

	// Save to cloudinary
	cld, _ := cloudinary.NewFromParams(os.Getenv("CLOUDINARY_NAME"), os.Getenv("CLOUDINARY_KEY"), os.Getenv("CLOUDINARY_SECRET"))
	result, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{PublicID: "colosach"})
	if err != nil {
		log.Fatal("cloudinary error: ", err.Error())
		return nil, err
	}
	img.Image_url = result.SecureURL

	// Download Image
	rawImage, err := utils.DownloadImage(img.Image_url)
	if err != nil {
		log.Fatal("error downloading image", err)
		return nil, err
	}

	// Generate color palette from the image
	colorsFromImage, err := prominentcolor.Kmeans(rawImage)
	if err != nil {
		log.Fatal("failed to generate colors: ", err)
		return nil, err
	}

	for _, c := range colorsFromImage {
		img.Palette = append(img.Palette, fmt.Sprintf("#%02x%02x%02x", c.Color.R, c.Color.G, c.Color.B))
	}
	// Save to mongo
	res, err := i.ImageCollection.InsertOne(i.ctx, img)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("image with that name already exists")
		}
		return nil, err
	}

	opt := options.Index()
	opt.SetUnique(true)

	index := mongo.IndexModel{Keys: bson.M{"name": 1}, Options: opt}
	if _, err := i.ImageCollection.Indexes().CreateOne(i.ctx, index); err != nil {
		return nil, errors.New("could not create index for name")
	}

	//todo: create index for tags too if we will query images by tags

	var newImage models.ImageResponse
	// newImage.Image_url = result.SecureURL
	query := bson.M{"_id": res.InsertedID}
	if err = i.ImageCollection.FindOne(i.ctx, query).Decode(&newImage); err != nil {
		return nil, err
	}
	return &newImage, nil
}
