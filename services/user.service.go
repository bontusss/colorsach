package services

import "github.com/bontusss/colosach/models"

type UserService interface {
	FindUserById(id string) (*models.DBResponse, error)
	FindUserByEmail(email string) (*models.DBResponse, error)
	UpdateUserById(id string, data *models.UpdateInput) (models.UserResponse, error)
	UpdateUserByEmail(email string, data *models.UpdateInput) (models.UserResponse, error)
}
