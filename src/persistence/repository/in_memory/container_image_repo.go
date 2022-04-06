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

func (repo *containerImageInMemoryRepository) Create(ctx context.Context, model model.ContainerImage) error {
	repo.fuzzedImages = append(repo.fuzzedImages, &model)
	return nil
}

func (repo *containerImageInMemoryRepository) Update(ctx context.Context, model model.ContainerImage) error {
	for i, savedImage := range repo.fuzzedImages {
		if savedImage.String() == model.String() {
			repo.fuzzedImages[i] = &model
			return nil
		}
	}

	return errors.New("couldn't find image to update")
}

func (repo *containerImageInMemoryRepository) FindByHash(ctx context.Context, hash string) (containerImage *model.ContainerImage, found bool, err error) {
	for _, image := range repo.fuzzedImages {
		if image.String() == hash {
			return image, true, nil
		}
	}

	return nil, false, nil
}
