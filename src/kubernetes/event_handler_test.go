package kubernetes

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/suecodelabs/cnfuzz/src/model"
	"github.com/suecodelabs/cnfuzz/src/persistence/repository"
	"github.com/suecodelabs/cnfuzz/src/persistence/repository/in_memory"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"testing"
)

func TestOnAdd(t *testing.T) {
	calledHandle := false
	handler := EventHandler{}
	handler.handleFunc = func(clientSet kubernetes.Interface, repositories *repository.Repositories, pod *apiv1.Pod) {
		calledHandle = true
	}
	handler.OnAdd(&apiv1.Pod{})
	assert.True(t, calledHandle)

}

func TestOnUpdate(t *testing.T) {
	calledHandle := false
	handler := EventHandler{}
	handler.handleFunc = func(clientSet kubernetes.Interface, repositories *repository.Repositories, pod *apiv1.Pod) {
		calledHandle = true
	}
	handler.OnUpdate(&apiv1.Pod{}, &apiv1.Pod{})
	assert.True(t, calledHandle)
}

func TestHandlePodEvent(t *testing.T) {
	// TODO
}

func TestFindPodInfoAndFuzz(t *testing.T) {
	// TODO
}

func TestContainsUnfuzzedImages(t *testing.T) {
	// func containsUnfuzzedImages(pod *apiv1.Pod, repo repository.IContainerImageRepository) (allImages []model.ContainerImage, containsUnfuzzedImages bool) {
	imageRepo := in_memory.CreateContainerImageRepository()
	existingImageName := "mycontainer"
	existingImage := model.ContainerImage{
		Hash:     "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		HashType: "sha256",
		Status:   model.Fuzzed,
		Tags:     []string{"latest"},
	}
	_, err := imageRepo.Create(context.TODO(), existingImage)
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
	allImages1, containsUnfuzzedImages1 := containsUnfuzzedImages(testPod1, imageRepo)
	assert.Equal(t, expectedResult1, containsUnfuzzedImages1)
	assert.Len(t, allImages1, len(testPod1.Status.ContainerStatuses))
	for _, image := range allImages1 {
		assert.Equal(t, model.BeingFuzzed, image.Status)
	}

	testPod2 := &apiv1.Pod{
		Status: apiv1.PodStatus{
			ContainerStatuses: []apiv1.ContainerStatus{
				{
					Image:   fmt.Sprintf("%s:%s", existingImageName, existingImage.Tags[0]),
					ImageID: fmt.Sprintf("docker-pullable://%s@%s:%s", existingImageName, existingImage.HashType, existingImage.Hash),
				},
			},
		},
	}
	// Get images currently in repo
	allImages2, containsUnfuzzedImages2 := containsUnfuzzedImages(testPod2, imageRepo)
	assert.Equal(t, false, containsUnfuzzedImages2)
	assert.Len(t, allImages2, len(testPod2.Status.ContainerStatuses))
	for _, image := range allImages2 {
		assert.Equal(t, model.Fuzzed, image.Status)
	}

}
