package repository

import (
	"context"
	"github.com/suecodelabs/cnfuzz/src/model"
	"github.com/suecodelabs/cnfuzz/src/persistence/repository/in_memory"
)

type BaseRepo[T any] interface {
	GetAll(ctx context.Context) ([]*T, error)
	Create(ctx context.Context, model T) (*T, error)
	Update(ctx context.Context, model T) (*T, error)
}

// IContainerImageRepository interface for a repository
type IContainerImageRepository interface {
	BaseRepo[model.ContainerImage]
	// FindByHash find a single container image by its name
	FindByHash(hash string) (containerImage *model.ContainerImage, found bool, err error)
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
