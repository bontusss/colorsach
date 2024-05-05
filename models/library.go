package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateLibraryRequest struct {
	Title       string           `json:"title" bson:"title" binding:"required"`
	Description string           `json:"description" bson:"description" binding:"required"`
	Owner       *UserResponse    `json:"owner" bson:"owner"`
	Image       *[]ImageResponse `json:"images,omitempty" bson:"images,omitempty"`
	Public      bool             `json:"public" bson:"public"`
	Likes       int              `json:"likes" bson:"likes"`
	Views       int              `json:"views" bson:"views"`
	Featured    bool             `json:"featured" bson:"featured"`
	CreatedAt   time.Time        `json:"created_at" bson:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt" bson:"updatedAt" `
}

type DBLibrary struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title" bson:"title" binding:"required"`
	Description string             `json:"description" bson:"description" binding:"required"`
	Owner       *UserResponse      `json:"owner" bson:"owner"`
	Images      *[]ImageResponse   `json:"images,omitempty" bson:"images,omitempty"`
	Public      bool               `json:"public" bson:"public"`
	Likes       int                `json:"likes" bson:"likes"`
	Views       int                `json:"views" bson:"views"`
	Featured    bool               `json:"featured" bson:"featured"`
	CreatedAt   time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type UpdateLibrary struct {
	Title       string           `json:"title" bson:"title" binding:"required"`
	Description string           `json:"description" bson:"description" binding:"required"`
	Images      *[]ImageResponse `json:"images,omitempty" bson:"images,omitempty"`
	Public      bool             `json:"public" bson:"public"`
	Featured    bool             `json:"featured" bson:"featured"`
	CreatedAt   time.Time        `json:"created_at" bson:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt" bson:"updatedAt"`
}
