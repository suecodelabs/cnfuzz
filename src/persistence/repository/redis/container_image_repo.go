package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/suecodelabs/cnfuzz/src/model"
	"strconv"
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
	key := containerImage.String()
	val := containerImage.Status
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

func (repo containerImageRedisRepository) FindByHash(ctx context.Context, key string) (containerImage *model.ContainerImage, found bool, err error) {
	val, err := repo.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	imgRepo := model.ContainerImageFromString(key)
	statusInt, convErr := strconv.ParseInt(val, 10, 16)
	if convErr != nil {
		return nil, true, convErr
	}
	imgRepo.Status = model.ImageFuzzStatus(statusInt)

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
