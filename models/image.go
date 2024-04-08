package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UploadImageInput struct {
	Name      string   `json:"name" bson:"name" binding:"required"`
	Tags      []string `json:"tags" bson:"tags"`
	Image_url string   `json:"image_url" bson:"image_url"`
	Palette   []string `json:"palette" bson:"palette"`
}

type ImageResponse struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Tags      []string           `json:"tags" bson:"tags"`
	Image_url string             `json:"image_url" bson:"image_url"`
}
