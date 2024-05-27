package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MetaData struct {
	DateTaken time.Time `json:"date_taken" bson:"date_taken"` // Date the image was taken
}

type Source struct {
	Original    string `json:"original" bson:"original"`
	Thumbnail   string `json:"thumbnail" bson:"thumbnail"`
	Watermarked string `json:"watermarked" bson:"watermarked"` // watermarked thumbnail
}

type UploadImageInput struct {
	Title       string     `json:"title" bson:"title" binding:"required"`
	Tags       []string   `json:"tags" bson:"tags"`
	Palette    []string   `json:"palette" bson:"palette"`
	Src        Source     `json:"source" bson:"source"`
	OwnerID    string     `json:"owner_id" bson:"owner_id"`
	Visibility Visibility `json:"visibility" bson:"visibility"`
	Published  bool       `json:"published" bson:"published"`
	UploadedAt time.Time  `json:"uploaded_at" bson:"uploaded_at"`
	UpdatedAt  time.Time  `json:"updated_at" bson:"updated_at"`
}

type ImageResponse struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	Title       string     `json:"title" bson:"title" binding:"required"`
	Tags       []string           `json:"tags" bson:"tags"`
	Likes      int                `json:"likes" bson:"likes"` // Number of times an image is liked
	OwnerID    string             `json:"owner_id" bson:"owner_id"`
	Palette    []string           `json:"palette" bson:"palette"`
	Src        Source     `json:"source" bson:"source"`
	Featured   bool               `json:"featured" bson:"featured"`
	Visibility Visibility         `json:"visibility" bson:"visibility"`
	MetaData   MetaData           `json:"meta_data" bson:"meta_data"`
	UploadedAt time.Time          `json:"uploaded_at" bson:"uploaded_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}
