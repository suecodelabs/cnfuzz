package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/suecodelabs/cnfuzz/src/model"
	"time"
)

type containerImageRedisRepository struct {
	client *redis.Client
}

func (repo containerImageRedisRepository) GetAll(ctx context.Context) ([]*model.ContainerImage, error) {
	//TODO implement me
	panic("function not implemented")
}

func (repo containerImageRedisRepository) Create(ctx context.Context, containerImage model.ContainerImage) error {
	key, val := containerImage.String()
	exp := time.Duration(0) // 0 means keep forever

	err := repo.client.Set(ctx, key, val, exp).Err()
	if err != nil {
		return err
	}
	return nil
}

func (repo containerImageRedisRepository) Update(ctx context.Context, containerImage model.ContainerImage) error {
	// This is lazy but updating and creating is the same function in redis ...
	return repo.Create(ctx, containerImage)
}

func (repo containerImageRedisRepository) FindByHash(ctx context.Context, hashKey string) (containerImage *model.ContainerImage, found bool, err error) {
	val, err := repo.client.Get(ctx, hashKey).Result()
	if err == redis.Nil {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	imgRepo, convErr := model.ContainerImageFromString(hashKey, val)
	if convErr != nil {
		return nil, true, convErr
	}

	return &imgRepo, true, nil
}

func CreateContainerImageRepository() *containerImageRedisRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &containerImageRedisRepository{
		client: rdb,
	}
}
