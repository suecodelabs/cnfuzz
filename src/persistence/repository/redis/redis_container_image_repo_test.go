package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"github.com/suecodelabs/cnfuzz/src/model"
	"strconv"
	"testing"
	"time"
)

var testContainerImage = model.ContainerImage{
	Hash:     "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
	HashType: "sha256",
	Status:   model.BeingFuzzed,
}

func createMockRepo() (client containerImageRedisRepository, mock redismock.ClientMock) {
	db, mock := redismock.NewClientMock()
	repo := containerImageRedisRepository{
		client: db,
	}

	return repo, mock
}

func TestCreateContainerImage(t *testing.T) {
	repo, mock := createMockRepo()

	expectedVal := strconv.Itoa(int(testContainerImage.Status))
	expectedKey := fmt.Sprintf("%s:%s", testContainerImage.HashType, testContainerImage.Hash)
	expectedExp := time.Duration(0)
	mock.ExpectSet(expectedKey, expectedVal, expectedExp).SetVal(expectedVal)
	err := repo.Create(context.TODO(), testContainerImage)
	assert.NoError(t, err)
}

func TestUpdateContainerImage(t *testing.T) {
	repo, mock := createMockRepo()

	expectedVal := strconv.Itoa(int(testContainerImage.Status))
	expectedKey := fmt.Sprintf("%s:%s", testContainerImage.HashType, testContainerImage.Hash)
	expectedExp := time.Duration(0)
	mock.ExpectSet(expectedKey, expectedVal, expectedExp).SetVal(expectedVal)
	err := repo.Update(context.TODO(), testContainerImage)
	assert.NoError(t, err)
}

func TestFindByHash(t *testing.T) {
	repo, mock := createMockRepo()

	strHash, strStatus := testContainerImage.String()

	// mock store returns an error if you don't add the ExpectSet line :(
	mock.ExpectSet(strHash, strStatus, time.Duration(0)).SetVal(strStatus)
	createErr := repo.client.Set(context.TODO(), strHash, strStatus, time.Duration(0)).Err()
	if !assert.NoError(t, createErr) {
		return
	}

	mock.ExpectGet(strHash).SetVal(strStatus)
	foundImage, didFind, findErr := repo.FindByHash(context.TODO(), strHash)

	assert.NoError(t, findErr)
	assert.True(t, didFind)
	assert.Equal(t, testContainerImage.Hash, foundImage.Hash)
	assert.Equal(t, testContainerImage.HashType, foundImage.HashType)
	assert.Equal(t, testContainerImage.Status, foundImage.Status)
}
