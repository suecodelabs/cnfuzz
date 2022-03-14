package repository

import (
	"github.com/suecodelabs/cnfuzz/src/model"
	"github.com/suecodelabs/cnfuzz/src/persistence/repository/in_memory"
)

type IContainerImageRepository interface {
	GetContainerImages() ([]model.ContainerImage, error)
	FindContainerImageByName(name string) (containerImage model.ContainerImage, found bool, err error)
	CreateContainerImage(image model.ContainerImage) error
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
