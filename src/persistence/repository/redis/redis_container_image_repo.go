// Copyright 2022 Sue B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/suecodelabs/cnfuzz/src/model"
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

func CreateContainerImageRepository(addr string, password string, db int) *containerImageRedisRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &containerImageRedisRepository{
		client: rdb,
	}
}
