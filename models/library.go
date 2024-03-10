package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateLibraryRequest struct {
	Title       string        `json:"title" bson:"title" binding:"required"`
	Description string        `json:"description" bson:"description" binding:"required"`
	Owner       *UserResponse `json:"owner" bson:"owner"`
	Image       string        `json:"image,omitempty" bson:"image,omitempty"`
	Public      bool          `json:"public" bson:"public"`
	CreatedAt   time.Time     `json:"created_at" bson:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt" bson:"updatedAt" `
}

type DBLibrary struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title" bson:"title" binding:"required"`
	Description string             `json:"description" bson:"description" binding:"required"`
	Owner       *UserResponse      `json:"owner" bson:"owner"`
	Image       string             `json:"image,omitempty" bson:"image,omitempty"`
	Public      bool               `json:"public" bson:"public"`
	CreatedAt   time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type UpdateLibrary struct {
	Title       string        `json:"title" bson:"title" binding:"required"`
	Description string        `json:"description" bson:"description" binding:"required"`
	Owner       *UserResponse `json:"owner" bson:"owner"`
	Image       string        `json:"image,omitempty" bson:"image,omitempty"`
	Public      bool          `json:"public" bson:"public"`
	CreatedAt   time.Time     `json:"created_at" bson:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt" bson:"updatedAt"`
}
