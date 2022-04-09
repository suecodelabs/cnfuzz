package in_memory

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/suecodelabs/cnfuzz/src/model"
	"testing"
)

func TestGetContainerImages(t *testing.T) {
	repo := createFilledMocRepo()
	returnedImages, err := repo.GetAll(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, len(repo.fuzzedImages), len(returnedImages))
	assert.Equal(t, repo.fuzzedImages, returnedImages)
}

func TestFindContainerImageByHashFound(t *testing.T) {
	img1, _ := model.CreateContainerImage("9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", "sha256", model.Fuzzed)
	img2, _ := model.CreateContainerImage("afa27b44d43b02a9fea41d13cedc2e4016cfcf87c5dbf990e593669aa8ce286d", "sha256", model.Fuzzed)
	images := []*model.ContainerImage{img1, img2}
	repo := createMocRepo(images)
	hashKey, _ := images[0].String()
	response, found, err := repo.FindByHash(context.TODO(), hashKey)
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, images[0], response)
}

func TestFindContainerImageByNameNil(t *testing.T) {
	var images []*model.ContainerImage
	repo := createMocRepo(images)
	response, found, err := repo.FindByHash(context.TODO(), "some-unknown-image-hash")
	assert.NoError(t, err)
	assert.False(t, found)
	assert.Empty(t, response)
}

func TestCreateContainerImage(t *testing.T) {
	repo := createFilledMocRepo()
	initLength := len(repo.fuzzedImages)
	newImage, _ := model.CreateContainerImage("9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", "sha256", model.BeingFuzzed)
	err := repo.Create(context.TODO(), *newImage)
	assert.NoError(t, err)
	assert.Greater(t, len(repo.fuzzedImages), initLength)
}

/* func TestCreateContainerImageIdFail(t *testing.T) {
	repo := createFilledMocRepo()
	initLength := len(repo.fuzzedImages)
	newImage := model.ContainerImage{
		Hash:     "",
		HashType: "sha256",
		Status:   model.Fuzzed,
		Tags:     []string{"latest"},
	}
	_, err := repo.Create(context.TODO(), newImage)
	assert.Error(t, err)
	assert.EqualError(t, err, "image hash is empty")
	assert.Equal(t, initLength, len(repo.fuzzedImages))
} */

/* func TestCreateContainerImageHashTypeFail(t *testing.T) {
	repo := createFilledMocRepo()
	initLength := len(repo.fuzzedImages)
	newImage := model.ContainerImage{
		Hash:     "afa27b44d43b02a9fea41d13cedc2e4016cfcf87c5dbf990e593669aa8ce286d",
		HashType: "",
		Status:   model.Fuzzed,
		Tags:     []string{"latest"},
	}
	_, err := repo.Create(context.TODO(), newImage)
	assert.Error(t, err)
	assert.EqualError(t, err, "image hash type can't be empty")
	assert.Equal(t, initLength, len(repo.fuzzedImages))
} */

func createMocRepo(containerImages []*model.ContainerImage) *containerImageInMemoryRepository {
	return &containerImageInMemoryRepository{
		fuzzedImages: containerImages,
	}
}

func createFilledMocRepo() *containerImageInMemoryRepository {
	img1, _ := model.CreateContainerImage("ac8f12f465a1300db7fbb2416bd922adc59a9c570ce8d54f8f7dd20ef2945464", "sha256", model.NotFuzzed)
	img2, _ := model.CreateContainerImage("3e27b58e2a4afe6db3020403403c1798adacb9adf0e60db2df27b120df521995", "sha256", model.Fuzzed)
	images := []*model.ContainerImage{img1, img2}

	return &containerImageInMemoryRepository{fuzzedImages: images}
}
