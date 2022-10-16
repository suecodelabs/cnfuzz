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

package in_memory

import (
	"context"
	"errors"
	"github.com/suecodelabs/cnfuzz/src/logger"
	"github.com/suecodelabs/cnfuzz/src/model"
)

type containerImageMem struct {
	l            logger.Logger
	fuzzedImages []*model.ContainerImage
}

func CreateContainerImageRepository(l logger.Logger) *containerImageMem {
	return &containerImageMem{
		l: l,
	}
}

func (repo *containerImageMem) GetAll(ctx context.Context) ([]*model.ContainerImage, error) {
	return repo.fuzzedImages, nil
}

func (repo *containerImageMem) Create(ctx context.Context, model model.ContainerImage) error {
	repo.fuzzedImages = append(repo.fuzzedImages, &model)
	return nil
}

func (repo *containerImageMem) Update(ctx context.Context, model model.ContainerImage) error {
	for i, savedImage := range repo.fuzzedImages {
		savedHash, _ := savedImage.String()
		targetHash, _ := model.String()
		if savedHash == targetHash {
			repo.fuzzedImages[i] = &model
			return nil
		}
	}

	return errors.New("couldn't find image to update")
}

func (repo *containerImageMem) GetByKey(ctx context.Context, key string) (containerImage *model.ContainerImage, found bool, err error) {
	for _, image := range repo.fuzzedImages {
		strHash, _ := image.String()
		if strHash == key {
			return image, true, nil
		}
	}

	return nil, false, nil
}
