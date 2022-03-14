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

func (repo containerImageInMemoryRepository) FindContainerImageByName(name string) (containerImage model.ContainerImage, found bool, err error) {
	for _, image := range repo.fuzzedImages {
		if image.Name == name {
			return image, true, nil
		}
	}

	return model.ContainerImage{}, false, nil
}

func (repo *containerImageInMemoryRepository) CreateContainerImage(image model.ContainerImage) error {
	// TODO: Move validation to model
	if len(image.Id) == 0 {
		return errors.New("image id is empty")
	}
	if len(image.Name) == 0 {
		return errors.New("image name is empty")
	}

	repo.fuzzedImages = append(repo.fuzzedImages, image)
	return nil
}

func (repo containerImageInMemoryRepository) UpdateContainerImage(image model.ContainerImage) error {
	for i, savedImage := range repo.fuzzedImages {
		if savedImage.Name == image.Name {
			repo.fuzzedImages[i] = image
			return nil
		}
	}

	return errors.New("couldn't find image to update")
}
