package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/bontusss/colosach/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestSignUpUser(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading app.env file", err)
	}

	// Create a new context and a new mongo client
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Create a new collection
	collection := client.Database("test").Collection("users")

	// Create a new AuthServiceImpl
	authService := NewAuthService(collection, ctx)

	// Create a new user
	user := &models.SignUpInput{
		Username: "testuser",
		Email:    "tEstuser@example.com",
		Password: "password",
	}

	// Call the SignUpUser method
	_, err = authService.SignUpUser(user)
	fmt.Println(user.Password)

	// Check if the error is as expected
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// Check if the user was inserted into the database
	var newUser *models.DBResponse
	query := bson.M{"email": user.Email}
	err = collection.FindOne(ctx, query).Decode(&newUser)
	if err != nil {
		t.Errorf("Expected to find user in database, but got %v", err)
	}

	// Check if the user's password was hashed
	if newUser.Password == "password" {
		t.Errorf("Expected password to be hashed, but it was not")
	}

	// Check if the user's email was set to lowercase
	if newUser.Email != user.Email {
		t.Errorf("Expected email to be lowercase, but it was not")
	}

	// Check if the user's role was set to UserRoleUser
	if newUser.Role != models.UserRoleUser {
		t.Errorf("Expected role to be UserRoleUser, but it was not")
	}

	// Check if the user's IsFirstLogin was set to true
	if !newUser.IsFirstLogin {
		t.Errorf("Expected IsFirstLogin to be true, but it was not")
	}
}
