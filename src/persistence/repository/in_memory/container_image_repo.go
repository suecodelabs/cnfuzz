package in_memory

import (
	"errors"

	"github.com/suecodelabs/cnfuzz/src/model"
)

type containerImageInMemoryRepository struct {
	fuzzedImages []model.ContainerImage
}

func CreateContainerImageRepository() *containerImageInMemoryRepository {

	return &containerImageInMemoryRepository{}
}

func (repo containerImageInMemoryRepository) GetContainerImages() ([]model.ContainerImage, error) {
	return repo.fuzzedImages, nil
}

func (repo containerImageInMemoryRepository) FindContainerImageByHash(hash string) (containerImage model.ContainerImage, found bool, err error) {
	for _, image := range repo.fuzzedImages {
		if image.Hash == hash {
			return image, true, nil
		}
	}

	return model.ContainerImage{}, false, nil
}

func (repo *containerImageInMemoryRepository) CreateContainerImage(image model.ContainerImage) error {
	valErr := image.Verify()
	if valErr != nil {
		return valErr
	}

	repo.fuzzedImages = append(repo.fuzzedImages, image)
	return nil
}

func (repo containerImageInMemoryRepository) UpdateContainerImage(image model.ContainerImage) error {
	for i, savedImage := range repo.fuzzedImages {
		if savedImage.Hash == image.Hash {
			repo.fuzzedImages[i] = image
			return nil
		}
	}

	return errors.New("couldn't find image to update")
}
