package services

import (
	"context"
	"errors"
	"fmt"
	"image"
	"log"
	"mime/multipart"
	"os"

	"github.com/Edlinorg/prominentcolor"
	"github.com/bontusss/colosach/models"
	"github.com/bontusss/colosach/utils"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/nfnt/resize"
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
func (i *ImageServiceImpl) UploadImage(img *models.UploadImageInput, file multipart.File, user *models.DBResponse) (*models.ImageResponse, error) {
	// Get image names from image title
	colorTagsFromImageName := utils.ExtractColorFromStrings(img.Title)
	img.Tags = append(img.Tags, colorTagsFromImageName...)

	uploadedImage, _, err := image.Decode(file)
	if err != nil {
		log.Println("unable to decode uploaded image")
		return nil, errors.New("unable to decode image")
	}

	// Create Image thumbnail from orginal image
	thumbnail := resize.Thumbnail(100, 100, uploadedImage, resize.Lanczos3)

	// Watermark thumbnail
	watermarkedThumbnail := utils.ApplyWatermark(thumbnail, fmt.Sprintf("by %s on Colosach", user.Username))

	file2, err := utils.ImageToMultipartFile(watermarkedThumbnail)
	if err != nil {
		log.Println("error converting thumbnail to a file")
		return nil, err
	}

	// Helper function to upload to Cloudinary
	uploadToCloudinary := func(file multipart.File, publicID string) (string, error) {
		cld, _ := cloudinary.NewFromParams(os.Getenv("CLOUDINARY_NAME"), os.Getenv("CLOUDINARY_KEY"), os.Getenv("CLOUDINARY_SECRET"))
		result, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{PublicID: publicID})
		if err != nil {
			log.Println("cloudinary error: ", err)
			return "", err
		}
		return result.SecureURL, nil
	}

	// Use the helper function for both original and thumbnail
	img.Src.Original, err = uploadToCloudinary(file, "colosach_original")
	if err != nil {
		return nil, err
	}
	img.Src.Thumbnail, err = uploadToCloudinary(file2, "colosach_thumbnail")
	if err != nil {
		return nil, err
	}

	// Generate color palette from the image
	colorsFromImage, err := prominentcolor.Kmeans(uploadedImage)
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

	index := mongo.IndexModel{Keys: bson.D{{Key: "title", Value: 1}, {Key: "owner_id", Value: 1}}, Options: opt}
	if _, err := i.ImageCollection.Indexes().CreateOne(i.ctx, index); err != nil {
		return nil, fmt.Errorf("could not create index for title: %v", err)
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
