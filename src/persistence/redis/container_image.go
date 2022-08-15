/*
 * Copyright 2022 Sue B.V.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package redis

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/suecodelabs/cnfuzz/src/health"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/suecodelabs/cnfuzz/src/model"
)

type containerImageRedis struct {
	l      logr.Logger
	client *redis.Client
}

func (repo containerImageRedis) Create(ctx context.Context, containerImage model.ContainerImage) error {
	key, val := containerImage.String()
	exp := time.Duration(0) // 0 means keep forever

	err := repo.client.Set(ctx, key, val, exp).Err()
	if err != nil {
		return err
	}
	return nil
}

func (repo containerImageRedis) Update(ctx context.Context, containerImage model.ContainerImage) error {
	// This is lazy but updating and creating is the same function in redis ...
	return repo.Create(ctx, containerImage)
}

func (repo containerImageRedis) GetByKey(ctx context.Context, key string) (obj *model.ContainerImage, found bool, err error) {
	val, err := repo.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	imgRepo, convErr := model.ContainerImageFromString(key, val)
	if convErr != nil {
		return nil, true, convErr
	}

	return &imgRepo, true, nil
}

func (repo containerImageRedis) CheckHealth(ctx context.Context) health.Health {
	status := repo.client.Ping(ctx)
	err := status.Err()
	if err != nil {
		h := health.NewHealth(false)
		h.Info[health.StatusKey] = health.UnHealthyStatus
		h.Info["reason"] = err.Error()
		return h
	} else {
		h := health.NewHealth(true)
		h.Info[health.StatusKey] = health.HealthyStatus
		return h
	}
}

func CreateContainerImageRedis(l logr.Logger, addr string, password string, db int) *containerImageRedis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &containerImageRedis{
		l:      l,
		client: rdb,
	}
}
