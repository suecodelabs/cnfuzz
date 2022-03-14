package in_memory

import (
	"github.com/stretchr/testify/assert"
	"github.com/suecodelabs/cnfuzz/src/model"
	"testing"
)

func TestGetContainerImages(t *testing.T) {
	repo := createFilledMocRepo()
	returnedImages, err := repo.GetContainerImages()
	assert.NoError(t, err)
	assert.Equal(t, len(repo.fuzzedImages), len(returnedImages))
	assert.Equal(t, repo.fuzzedImages, returnedImages)
}

func TestFindContainerImageByNameFound(t *testing.T) {
	targetName := "someimage"
	images := []model.ContainerImage{
		{
			Id:       targetName + "-123456789",
			Name:     targetName,
			Versions: []model.ContainerImageVersion{{}},
		},
		{
			Id:       "someotherimage-123456789",
			Name:     "someotherimage",
			Versions: []model.ContainerImageVersion{{}},
		},
	}
	repo := createMocRepo(images)
	response, found, err := repo.FindContainerImageByName(targetName)
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, images[0], response)
}

func TestFindContainerImageByNameNil(t *testing.T) {
	var images []model.ContainerImage
	repo := createMocRepo(images)
	response, found, err := repo.FindContainerImageByName("some-unknown-image")
	assert.NoError(t, err)
	assert.False(t, found)
	assert.Empty(t, response.Id)
	assert.Empty(t, response.Name)
}

func TestCreateContainerImage(t *testing.T) {
	repo := createFilledMocRepo()
	initLength := len(repo.fuzzedImages)
	newImage := model.ContainerImage{
		Id:       "somenewimage-123456789",
		Name:     "somenewimage",
		Versions: []model.ContainerImageVersion{{}},
	}
	err := repo.CreateContainerImage(newImage)
	assert.NoError(t, err)
	assert.Greater(t, len(repo.fuzzedImages), initLength)
}

func TestCreateContainerImageIdFail(t *testing.T) {
	repo := createFilledMocRepo()
	initLength := len(repo.fuzzedImages)
	newImage := model.ContainerImage{
		Id:       "",
		Name:     "someimage",
		Versions: []model.ContainerImageVersion{{}},
	}
	err := repo.CreateContainerImage(newImage)
	assert.Error(t, err)
	assert.EqualError(t, err, "image id is empty")
	assert.Equal(t, initLength, len(repo.fuzzedImages))
}

func TestCreateContainerImageNameFail(t *testing.T) {
	repo := createFilledMocRepo()
	initLength := len(repo.fuzzedImages)
	newImage := model.ContainerImage{
		Id:       "someimage-123456789",
		Name:     "",
		Versions: []model.ContainerImageVersion{{}},
	}
	err := repo.CreateContainerImage(newImage)
	assert.Error(t, err)
	assert.EqualError(t, err, "image name is empty")
	assert.Equal(t, initLength, len(repo.fuzzedImages))
}

func createMocRepo(containerImages []model.ContainerImage) *containerImageInMemoryRepository {
	return &containerImageInMemoryRepository{
		fuzzedImages: containerImages,
	}
}

func createFilledMocRepo() *containerImageInMemoryRepository {
	images := []model.ContainerImage{
		{
			Id:       "someimage-123456789",
			Name:     "someimage",
			Versions: []model.ContainerImageVersion{{}},
		},
		{
			Id:       "someotherimage-123456789",
			Name:     "someotherimage",
			Versions: []model.ContainerImageVersion{{}},
		},
	}
	return &containerImageInMemoryRepository{fuzzedImages: images}
}
