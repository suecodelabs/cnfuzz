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

package persistence

import (
	"context"
	"fmt"
	"github.com/suecodelabs/cnfuzz/src/pkg/cnfuzz/health"
	"github.com/suecodelabs/cnfuzz/src/pkg/cnfuzz/model"
	"github.com/suecodelabs/cnfuzz/src/pkg/cnfuzz/persistence/in_memory"
	"github.com/suecodelabs/cnfuzz/src/pkg/cnfuzz/persistence/redis"
	"github.com/suecodelabs/cnfuzz/src/pkg/logger"
)

type Cache[T any] interface {
	Create(ctx context.Context, model T) error
	Update(ctx context.Context, model T) error
	GetByKey(ctx context.Context, key string) (obj *T, found bool, err error)
}

type Storage struct {
	ContainerImageCache Cache[model.ContainerImage]
}

func InitRedisCache(l logger.Logger, addr string, port string, hc health.Checker) *Storage {
	if len(port) > 0 {
		addr = fmt.Sprintf("%s:%s", addr, port)
	}
	pass := ""
	db := 0
	l.V(logger.DebugLevel).Info(fmt.Sprintf("using redis from %s", addr), "redisAddr", addr, "dbId", db)
	cICache := redis.CreateContainerImageRedis(l, addr, pass, db)

	hc.RegisterCheck("redis", cICache)
	return &Storage{ContainerImageCache: cICache}
}

func InitMemoryCache(l logger.Logger) *Storage {
	cICache := in_memory.CreateContainerImageRepository(l)
	return &Storage{ContainerImageCache: cICache}
}
