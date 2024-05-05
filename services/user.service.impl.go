package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/bontusss/colosach/models"
	"github.com/bontusss/colosach/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserServiceImpl(collection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{collection, ctx}
}

func (uc *UserServiceImpl) FindUserById(id string) (*models.DBResponse, error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var user *models.DBResponse

	query := bson.M{"_id": oid}
	err := uc.collection.FindOne(uc.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

// UpdateUserById implements UserService.
func (uc *UserServiceImpl) UpdateUserById(id string, data *models.UpdateInput) (models.UserResponse, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.UserResponse{}, fmt.Errorf("invalid ID format: %v", err)
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	updateFields, err := utils.ToDoc(data)
	if err != nil {
		return models.UserResponse{}, fmt.Errorf("error converting data to BSON document: %v", err)
	}

	if updateFields == nil || len(*updateFields) == 0 {
		return models.UserResponse{}, fmt.Errorf("no data provided to update")
	}

	update := bson.M{"$set": updateFields}
	result, err := uc.collection.UpdateOne(uc.ctx, filter, update)
	if err != nil {
		return models.UserResponse{}, fmt.Errorf("error in UpdateOne: %v", err)
	}

	if result.ModifiedCount == 0 {
		return models.UserResponse{}, fmt.Errorf("no user found with the given ID or no update needed")
	}

	var updatedUser models.DBResponse
	err = uc.collection.FindOne(uc.ctx, filter).Decode(&updatedUser)
	if err != nil {
		return models.UserResponse{}, fmt.Errorf("error in FindOne: %v", err)
	}

	return models.FilteredResponse(&updatedUser), nil
}

func (uc *UserServiceImpl) FindUserByEmail(email string) (*models.DBResponse, error) {
	var user *models.DBResponse

	query := bson.M{"email": strings.ToLower(email)}
	err := uc.collection.FindOne(uc.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

func (uc *UserServiceImpl) UpdateUserByEmail(email string, data *models.UpdateInput) (models.UserResponse, error) {
	// Create an update query
	filter := bson.D{{Key: "email", Value: email}}
	update := bson.M{"$set": bson.M{"role": data.Role, "updated_at": data.UpdatedAt}}

	// Execute the update operation
	result, err := uc.collection.UpdateOne(uc.ctx, filter, update)
	if err != nil {
		return models.UserResponse{}, fmt.Errorf("error in UpdateOne: %v", err)
	}

	// Check if the document was successfully updated
	if result.ModifiedCount == 0 {
		return models.UserResponse{}, fmt.Errorf("no user found with the given email or no update needed")
	}

	var updatedUser models.DBResponse
	err = uc.collection.FindOne(uc.ctx, filter).Decode(&updatedUser)
	if err != nil {
		return models.UserResponse{}, fmt.Errorf("error in FindOne: %v", err)
	}
	return models.FilteredResponse(&updatedUser), nil
}
