package services

import "github.com/bontusss/colosach/models"

type LibraryService interface {
	CreateLibrary(*models.CreateLibraryRequest) (*models.DBLibrary, error)
	UpdateLibrary(library *models.UpdateLibrary) (*models.DBLibrary, error)
	FindLibraryByID(string) (*models.DBLibrary, error)
	FindLibraries(page int, limit int) ([]*models.DBLibrary, error)
	DeleteLibrary(string) error
}
