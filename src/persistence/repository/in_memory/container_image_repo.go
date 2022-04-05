package in_memory

import (
	"context"
	"errors"

	"github.com/suecodelabs/cnfuzz/src/model"
)

type containerImageInMemoryRepository struct {
	fuzzedImages []*model.ContainerImage
}

func CreateContainerImageRepository() *containerImageInMemoryRepository {

	return &containerImageInMemoryRepository{}
}

func (repo *containerImageInMemoryRepository) GetAll(ctx context.Context) ([]*model.ContainerImage, error) {
	return repo.fuzzedImages, nil
}

func (repo *containerImageInMemoryRepository) Create(ctx context.Context, model model.ContainerImage) (*model.ContainerImage, error) {
	image := &model
	repo.fuzzedImages = append(repo.fuzzedImages, &model)
	return image, nil
}

func (repo *containerImageInMemoryRepository) Update(ctx context.Context, model model.ContainerImage) (*model.ContainerImage, error) {
	for i, savedImage := range repo.fuzzedImages {
		if savedImage.Hash == model.Hash {
			repo.fuzzedImages[i] = &model
			return repo.fuzzedImages[i], nil
		}
	}

	return nil, errors.New("couldn't find image to update")
}

func (repo *containerImageInMemoryRepository) FindByHash(hash string) (containerImage *model.ContainerImage, found bool, err error) {
	for _, image := range repo.fuzzedImages {
		if image.Hash == hash {
			return image, true, nil
		}
	}

	return nil, false, nil
}
