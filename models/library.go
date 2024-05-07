package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Visibility string

const (
	IsPublic  Visibility = "public"
	IsPrivate Visibility = "private"
)

type CreateLibraryRequest struct {
	Name        string     `json:"name" bson:"name" binding:"required"`
	Description string     `json:"description" bson:"description" binding:"required"`
	OwnerID     string     `json:"owner_id" bson:"owner_id"`
	ImageIDs    []string   `json:"images,omitempty" bson:"images,omitempty"`
	Visibility  Visibility `json:"visibility" bson:"visibility"`
	Likes       int        `json:"likes" bson:"likes"`
	Views       int        `json:"views" bson:"views"`
	Featured    bool       `json:"featured" bson:"featured"`
	CreatedAt   time.Time  `json:"created_at" bson:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt" bson:"updatedAt" `
}

type DBLibrary struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name" binding:"required"`
	Description string             `json:"description" bson:"description" binding:"required"`
	OwnerID     string             `json:"owner_id" bson:"owner_id"`
	ImageIDs    []string           `json:"images,omitempty" bson:"images,omitempty"`
	Visibility  Visibility         `json:"visibility" bson:"visibility"`
	Likes       int                `json:"likes" bson:"likes"`
	Views       int                `json:"views" bson:"views"`
	Featured    bool               `json:"featured" bson:"featured"`
	CreatedAt   time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type UpdateLibrary struct {
	Name        string     `json:"name" bson:"name" binding:"required"`
	Description string     `json:"description" bson:"description" binding:"required"`
	ImageIDs    []string   `json:"images,omitempty" bson:"images,omitempty"`
	Visibility  Visibility `json:"visibility" bson:"visibility"`
	Featured    bool       `json:"featured" bson:"featured"`
	CreatedAt   time.Time  `json:"created_at" bson:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt" bson:"updatedAt"`
}
