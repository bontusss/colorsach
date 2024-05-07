package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bontusss/colosach/models"
	"github.com/bontusss/colosach/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LibraryServiceImpl struct {
	collection *mongo.Collection
	ctx               context.Context
}

func NewLibService(libCollection *mongo.Collection, ctx context.Context) LibraryService {
	return &LibraryServiceImpl{libCollection, ctx}
}

func (ls *LibraryServiceImpl) UpdateLibrary(id string, library *models.UpdateLibrary) (*models.DBLibrary, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", err)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	updateFields, err := utils.ToDoc(library)
	if err != nil {
		return nil, fmt.Errorf("error converting data to bson document: %v", err)
	}

	if updateFields == nil || len(*updateFields) == 0 {
		return nil, errors.New("no data provided to update")
	}

	update := bson.M{"$set": updateFields}
	result, err := ls.collection.UpdateOne(ls.ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("error updating library: %v", err)
	}

	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("no library with that id")
	}

	var updatedLibrary models.DBLibrary
	err = ls.collection.FindOne(ls.ctx, filter).Decode(&updatedLibrary)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	return &updatedLibrary, nil
}

func (ls *LibraryServiceImpl) FindLibraryByID(id string) (*models.DBLibrary, error) {
	var library *models.DBLibrary
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": oid}
	err = ls.collection.FindOne(ls.ctx, filter).Decode(&library)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}
	return library, nil
}

func (ls *LibraryServiceImpl) FindLibraries(page int, limit int) ([]*models.DBLibrary, error) {
	//TODO implement me
	panic("implement me")
}

func (ls *LibraryServiceImpl) DeleteLibrary(s string) error {
	//TODO implement me
	panic("implement me")
}

func (ls *LibraryServiceImpl) CreateLibrary(lib *models.CreateLibraryRequest) (*models.DBLibrary, error) {
	lib.CreatedAt = time.Now()
	lib.UpdatedAt = lib.CreatedAt
	lib.Visibility = models.IsPublic // visibility defaults to public if not set
	lib.Featured = false

	opt := options.Index()
	opt.SetUnique(true)

	index := mongo.IndexModel{Keys: bson.D{{Key: "name", Value: 1}, {Key: "owner_id", Value: 1}}, Options: opt}
	if _, err := ls.collection.Indexes().CreateOne(ls.ctx, index); err != nil {
		return nil, fmt.Errorf("could not create index for name: %v", err)
	}

	res, err := ls.collection.InsertOne(ls.ctx, lib)
	if err != nil {
		if mongoErr, ok := err.(mongo.WriteException); ok {
			for _, we := range mongoErr.WriteErrors {
				if we.Code == 11000 { // Duplicate key error code
					return nil, fmt.Errorf("you already have a library with that name")
				}
			}
		}
		return nil, fmt.Errorf("error inserting new library: %v", err)
	}


	var newLibrary *models.DBLibrary
	query := bson.M{"_id": res.InsertedID}
	if err = ls.collection.FindOne(ls.ctx, query).Decode(&newLibrary); err != nil {
		return nil, err
	}

	return newLibrary, nil
}
