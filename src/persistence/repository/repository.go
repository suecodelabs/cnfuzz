package repository

import (
	"github.com/suecodelabs/cnfuzz/src/model"
	"github.com/suecodelabs/cnfuzz/src/persistence/repository/in_memory"
)

// IContainerImageRepository interface for a repository
type IContainerImageRepository interface {
	// GetContainerImages get all the container images
	GetContainerImages() ([]model.ContainerImage, error)
	// FindContainerImageByName find a single container image by its name
	FindContainerImageByName(name string) (containerImage model.ContainerImage, found bool, err error)
	// CreateContainerImage create a new container image
	CreateContainerImage(image model.ContainerImage) error
	// UpdateContainerImage edit an existing container image
	UpdateContainerImage(image model.ContainerImage) error
}

// Repositories contains all the repo structs
type Repositories struct {
	ContainerImageRepo IContainerImageRepository
}

// InitRepositories should be called only ones
func InitRepositories() *Repositories {
	containerImageRepo := in_memory.CreateContainerImageRepository()
	return &Repositories{ContainerImageRepo: containerImageRepo}
}
