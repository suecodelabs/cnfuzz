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

package k8s

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/suecodelabs/cnfuzz/src/internal/model"
	"github.com/suecodelabs/cnfuzz/src/internal/persistence"
	"github.com/suecodelabs/cnfuzz/src/internal/persistence/in_memory"
	"github.com/suecodelabs/cnfuzz/src/pkg/config"
	"github.com/suecodelabs/cnfuzz/src/pkg/logger"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"testing"
)

func TestOnAdd(t *testing.T) {
	calledHandle := false
	c := controller{}
	c.handleFunc = func(l logger.Logger, clientSet kubernetes.Interface, storage *persistence.Storage, config *config.CnFuzzConfig, overwrites config.Overwrites, pod *apiv1.Pod) {
		calledHandle = true
	}
	c.OnAdd(&apiv1.Pod{})
	assert.True(t, calledHandle)

}

func TestOnUpdate(t *testing.T) {
	calledHandle := false
	c := controller{}
	c.handleFunc = func(l logger.Logger, clientSet kubernetes.Interface, storage *persistence.Storage, config *config.CnFuzzConfig, overwrites config.Overwrites, pod *apiv1.Pod) {
		calledHandle = true
	}
	c.OnUpdate(&apiv1.Pod{}, &apiv1.Pod{})
	assert.True(t, calledHandle)
}

func TestHandlePodEvent(t *testing.T) {
	// TODO
}

func TestFindPodInfoAndFuzz(t *testing.T) {
	// TODO
}

func TestContainsUnfuzzedImages(t *testing.T) {
	l := logger.CreateDebugLogger()
	// func containsUnfuzzedImages(pod *apiv1.Pod, repo repository.IContainerImageRepository) (allImages []model.ContainerImage, containsUnfuzzedImages bool) {
	imageRepo := in_memory.CreateContainerImageRepository(l)
	existingImageName := "mycontainer"
	existingImage, _ := model.CreateContainerImage("9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", "sha256", model.Fuzzed)
	err := imageRepo.Create(context.TODO(), existingImage)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to prep image repo for ContainsUnfuzzedImages function test: %w", err))
	}

	expectedResult1 := true
	testPod1 := &apiv1.Pod{
		Status: apiv1.PodStatus{
			ContainerStatuses: []apiv1.ContainerStatus{
				{
					Image:   "myregistry/apiimage:latest:debian",
					ImageID: "docker-pullable://myregistry/apiimage@sha256:729610843b7af92d6c481af4e066cb3d4dfabbe8de7d29f58e8cff2f7170115b",
				},
				{
					Image:   "myregistry/dbimage:latest",
					ImageID: "docker-pullable://myregistry/dbimage@sha256:64ebf2c8187f48e2d919653e9c43c830c7b2cd6418e5ad815108dfe79863a94",
				},
			},
		},
	}
	allImages1, containsUnfuzzedImages1 := containsUnfuzzedImages(l, testPod1, imageRepo)
	assert.Equal(t, expectedResult1, containsUnfuzzedImages1)
	assert.Len(t, allImages1, len(testPod1.Status.ContainerStatuses))
	for _, image := range allImages1 {
		assert.Equal(t, model.BeingFuzzed, image.Status)
	}

	testPod2 := &apiv1.Pod{
		Status: apiv1.PodStatus{
			ContainerStatuses: []apiv1.ContainerStatus{
				{
					Image:   fmt.Sprintf("%s:%s", existingImageName, "latest"),
					ImageID: fmt.Sprintf("docker-pullable://%s@%s:%s", existingImageName, existingImage.HashType, existingImage.Hash),
				},
			},
		},
	}
	// Get images currently in repo
	allImages2, containsUnfuzzedImages2 := containsUnfuzzedImages(l, testPod2, imageRepo)
	assert.Equal(t, false, containsUnfuzzedImages2)
	assert.Len(t, allImages2, len(testPod2.Status.ContainerStatuses))
	for _, image := range allImages2 {
		assert.Equal(t, model.Fuzzed, image.Status)
	}

}
