package model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
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
	for _, image := range result {
		// Container can only have one image
		if assert.Len(t, image.Versions, 1) {
			assert.Equal(t, image.Versions[0].Status, Unknown)
		}

	}
}
