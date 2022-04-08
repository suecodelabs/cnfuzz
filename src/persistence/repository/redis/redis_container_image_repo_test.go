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

func TestCreateContainerImage(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := containerImageRedisRepository{
		client: db,
	}

	newImage := model.ContainerImage{
		Hash:     "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		HashType: "sha256",
		Status:   model.BeingFuzzed,
	}

	expectedVal := strconv.Itoa(int(newImage.Status))
	expectedKey := fmt.Sprintf("%s:%s", newImage.HashType, newImage.Hash)
	expectedExp := time.Duration(0)
	mock.ExpectSet(expectedKey, expectedVal, expectedExp).SetVal(expectedVal)
	err := repo.Create(context.TODO(), newImage)
	assert.NoError(t, err)
}

func TestFindByHash(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := containerImageRedisRepository{
		client: db,
	}

	newImage := model.ContainerImage{
		Hash:     "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		HashType: "sha256",
		Status:   model.BeingFuzzed,
	}

	strHash, strStatus := newImage.String()

	// mock store returns an error if you don't add the ExpectSet line :(
	mock.ExpectSet(strHash, strStatus, time.Duration(0)).SetVal(strStatus)
	createErr := db.Set(context.TODO(), strHash, strStatus, time.Duration(0)).Err()
	if !assert.NoError(t, createErr) {
		return
	}

	mock.ExpectGet(strHash).SetVal(strStatus)
	foundImage, didFind, findErr := repo.FindByHash(context.TODO(), strHash)

	assert.NoError(t, findErr)
	assert.True(t, didFind)
	assert.Equal(t, newImage.Hash, foundImage.Hash)
	assert.Equal(t, newImage.HashType, foundImage.HashType)
	assert.Equal(t, newImage.Status, foundImage.Status)
}
