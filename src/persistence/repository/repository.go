package repository

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"github.com/suecodelabs/cnfuzz/src/cmd"
	"github.com/suecodelabs/cnfuzz/src/model"
	"github.com/suecodelabs/cnfuzz/src/persistence/repository/in_memory"
	"github.com/suecodelabs/cnfuzz/src/persistence/repository/redis"
)

type BaseRepo[T any] interface {
	GetAll(ctx context.Context) ([]*T, error)
	Create(ctx context.Context, model T) error
	Update(ctx context.Context, model T) error
}

// IContainerImageRepository interface for a ContainerImage repository that uses the BaseRepo
type IContainerImageRepository interface {
	BaseRepo[model.ContainerImage]
	// FindByHash find a single container image by its hash key (format: hashtype:hash)
	FindByHash(ctx context.Context, hashKey string) (containerImage *model.ContainerImage, found bool, err error)
}

// Repositories contains all the repo structs
type Repositories struct {
	ContainerImageRepo IContainerImageRepository
}

// InitRepositories should be called only ones
func InitRepositories(repoType RepoType) *Repositories {
	if repoType == Redis {
		// Would prefer to get the config in some other way
		addr := viper.GetString(cmd.RedisHostName)
		port := viper.GetString(cmd.RedisPort)
		if len(port) > 0 {
			addr = fmt.Sprintf("%s:%s", addr, port)
		}
		pass := ""
		db := 0
		containerImageRepo := redis.CreateContainerImageRepository(addr, pass, db)
		return &Repositories{ContainerImageRepo: containerImageRepo}
	} else if repoType == InMemory {
		containerImageRepo := in_memory.CreateContainerImageRepository()
		return &Repositories{ContainerImageRepo: containerImageRepo}
	}
	return nil
}
