// Copyright 2022 Sue B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	"strconv"
	"testing"
)

func TestCreateContainerImagesFromPod(t *testing.T) {
	testRegistry := "myregistry"
	hashType := "sha256"
	testPod := &apiv1.Pod{
		Status: apiv1.PodStatus{
			ContainerStatuses: []apiv1.ContainerStatus{
				{
					Image:   fmt.Sprintf("%s/apiimage:latest:debian", testRegistry),
					ImageID: fmt.Sprintf("docker-pullable://%s/apiimage@%s:729610843b7af92d6c481af4e066cb3d4dfabbe8de7d29f58e8cff2f7170115b", testRegistry, hashType),
				},
				{
					Image:   fmt.Sprintf("%s/dbimage:latest", testRegistry),
					ImageID: fmt.Sprintf("docker-pullable://%s/dbimage@%s:64ebf2c8187f48e2d919653e9c43c830c7b2cd6418e5ad815108dfe79863a94", testRegistry, hashType),
				},
			},
		},
	}
	result, err := CreateContainerImagesFromPod(testPod)
	assert.NoError(t, err)
	assert.Len(t, result, len(testPod.Status.ContainerStatuses))
	assert.Equal(t, result[0].Status, NotFuzzed)
}

func TestContainerImage_Verify(t *testing.T) {
	newImage := ContainerImage{
		Hash:     "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		HashType: "sha256",
		Status:   Fuzzed,
	}
	err := newImage.Verify()
	assert.NoError(t, err, "this container image is valid, but the verify function says differently")
}

func TestContainerImage_VerifyIdError(t *testing.T) {
	newImage := ContainerImage{
		Hash:     "",
		HashType: "sha256",
		Status:   Fuzzed,
	}
	err := newImage.Verify()
	assert.Error(t, err)
	assert.EqualError(t, err, "container image is invalid because image hash is empty")
}

func TestContainerImage_VerifyHashTypeError(t *testing.T) {
	newImage := ContainerImage{
		Hash:     "afa27b44d43b02a9fea41d13cedc2e4016cfcf87c5dbf990e593669aa8ce286d",
		HashType: "",
		Status:   Fuzzed,
	}
	err := newImage.Verify()
	assert.Error(t, err)
	assert.EqualError(t, err, "container image is invalid because image hash type can't be empty")
}

func TestString(t *testing.T) {
	testImage := ContainerImage{
		Hash:     "afa27b44d43b02a9fea41d13cedc2e4016cfcf87c5dbf990e593669aa8ce286d",
		HashType: "sha256",
		Status:   NotFuzzed,
	}

	strHash, strStatus := testImage.String()
	expectedHash := fmt.Sprintf("%s:%s", testImage.HashType, testImage.Hash)
	expectedStatus := strconv.Itoa(int(testImage.Status))
	assert.Equal(t, expectedHash, strHash, "string method returns unexpected format or invalid values")
	assert.Equal(t, expectedStatus, strStatus, "string method returns unexpected status")
}

func TestContainerImageFromString(t *testing.T) {
	hash := "afa27b44d43b02a9fea41d13cedc2e4016cfcf87c5dbf990e593669aa8ce286d"
	hashType := "sha256"
	status := ImageFuzzStatus(1)
	imageStr := fmt.Sprintf("%s:%s", hashType, hash)
	createdImg, convErr := ContainerImageFromString(imageStr, strconv.Itoa(int(status)))
	if assert.NoError(t, convErr) {
		assert.Equal(t, hashType, createdImg.HashType)
		assert.Equal(t, hash, createdImg.Hash)
		assert.Equal(t, status, createdImg.Status)
	}
}
