package services

import (
	"context"
	"errors"
	"github.com/bontusss/colosach/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type LibraryServiceImpl struct {
	LibraryCollection *mongo.Collection
	ctx               context.Context
}

func NewLibService(libCollection *mongo.Collection, ctx context.Context) LibraryService {
	return &LibraryServiceImpl{libCollection, ctx}
}

func (ls *LibraryServiceImpl) UpdateLibrary(library *models.UpdateLibrary) (*models.DBLibrary, error) {
	//TODO implement me
	panic("implement me")
}

func (ls *LibraryServiceImpl) FindLibraryByID(s string) (*models.DBLibrary, error) {
	//TODO implement me
	panic("implement me")
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
	lib.Public = true

	res, err := ls.LibraryCollection.InsertOne(ls.ctx, lib)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("library with that title already exists")
		}
		return nil, err
	}

	opt := options.Index()
	opt.SetUnique(true)

	index := mongo.IndexModel{Keys: bson.M{"title": 1}, Options: opt}
	if _, err := ls.LibraryCollection.Indexes().CreateOne(ls.ctx, index); err != nil {
		return nil, errors.New("could not create index for title")
	}

	var newLibrary *models.DBLibrary
	query := bson.M{"_id": res.InsertedID}
	if err = ls.LibraryCollection.FindOne(ls.ctx, query).Decode(&newLibrary); err != nil {
		return nil, err
	}

	return newLibrary, nil
}
