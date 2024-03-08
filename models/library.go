package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CreateLibraryRequest struct {
	Title       string      `json:"title" bson:"title" binding:"required"`
	Description string      `json:"description" bson:"description" binding:"required"`
	User        *DBResponse `json:"user" bson:"user" binding:"required"`
	Image       string      `json:"image,omitempty" bson:"image,omitempty"`
	Public      bool        `json:"public" bson:"public" binding:"required"`
	CreatedAt   time.Time   `json:"created_at" bson:"createdAt" binding:"required"`
	UpdatedAt   time.Time   `json:"updatedAt" bson:"updatedAt" binding:"required"`
}

type DBLibrary struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title" bson:"title" binding:"required"`
	Description string             `json:"description" bson:"description" binding:"required"`
	User        *DBResponse        `json:"user" bson:"user" binding:"required"`
	Image       string             `json:"image,omitempty" bson:"image,omitempty"`
	Public      bool               `json:"public" bson:"public" binding:"required"`
	CreatedAt   time.Time          `json:"created_at" bson:"createdAt" binding:"required"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt" binding:"required"`
}

type UpdateLibrary struct {
	Title       string      `json:"title" bson:"title" binding:"required"`
	Description string      `json:"description" bson:"description" binding:"required"`
	User        *DBResponse `json:"user" bson:"user" binding:"required"`
	Image       string      `json:"image,omitempty" bson:"image,omitempty"`
	Public      bool        `json:"public" bson:"public" binding:"required"`
	CreatedAt   time.Time   `json:"created_at" bson:"createdAt" binding:"required"`
	UpdatedAt   time.Time   `json:"updatedAt" bson:"updatedAt" binding:"required"`
}
