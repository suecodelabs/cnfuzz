package model

import (
	"errors"
	"fmt"
	kutil "github.com/suecodelabs/cnfuzz/src/kubernetes/util"
	"github.com/suecodelabs/cnfuzz/src/log"
	apiv1 "k8s.io/api/core/v1"
	"strconv"
	"strings"
)

type ImageFuzzStatus int

const (
	NotFuzzed ImageFuzzStatus = iota
	Fuzzed
	BeingFuzzed
)

// ContainerImage container image struct that contains a fuzz status
type ContainerImage struct {
	Hash     string
	HashType string
	Status   ImageFuzzStatus
}

// Verify verifies if this model is valid, model is not valid if some non nullable properties are empty or nil
func (img ContainerImage) Verify() error {
	if len(img.Hash) == 0 {
		return errors.New("container image is invalid because image hash is empty")
	}
	if len(img.HashType) == 0 {
		return errors.New("container image is invalid because image hash type can't be empty")
	}
	return nil
}

// String ContainerImage to string representation
func (img ContainerImage) String() (key string, status string) {
	return fmt.Sprintf("%s:%s", img.HashType, img.Hash), strconv.Itoa(int(img.Status))
}

// ContainerImageFromString create ContainerImage from a string in format hashtype:hash and a string representing the status as an int
func ContainerImageFromString(hashString string, statusString string) (image ContainerImage, convErr error) {
	hashSplit := strings.Split(hashString, ":")
	hashType := hashSplit[0]
	hash := hashSplit[1]
	status, convErr := strconv.ParseInt(statusString, 10, 16)
	if convErr != nil {
		return ContainerImage{}, convErr
	}
	return ContainerImage{
		Hash:     hash,
		HashType: hashType,
		Status:   ImageFuzzStatus(status),
	}, nil
}

func CreateContainerImage(hash string, hashType string, status ImageFuzzStatus) (ContainerImage, error) {
	img := ContainerImage{
		Hash:     hash,
		HashType: hashType,
		Status:   status,
	}
	return img, img.Verify()
}

// CreateContainerImagesFromPod extracts container info from a pod and converts their images to ContainerImages
func CreateContainerImagesFromPod(pod *apiv1.Pod) ([]ContainerImage, error) {
	logger := log.L()
	var images []ContainerImage
mainloop:
	for _, status := range pod.Status.ContainerStatuses {
		if len(status.ImageID) == 0 || len(status.Image) == 0 {
			// TODO are there other image names or ID's that we can't parse/check?
			logger.Warnf("image ID \"%s\" or image name \"%s\" are invalid and can't be checked", status.ImageID, status.Image)
			logger.Warnf("if this is the only container in the pod, the fuzzer won't fuzz this pod")
			continue
		}

		hash, hashType := kutil.SplitImageId(status.ImageID)

		// Look for duplicate image hashes/versions
		for _, image := range images {
			if image.Hash == hash {
				// Image already exists in the image array, so avoid creating a duplicate
				continue mainloop
			}
		}

		newImage, createErr := CreateContainerImage(hash, hashType, NotFuzzed)
		if createErr != nil {
			// This means either the hash or hash type is empty
			// Returning the error instead of ignoring it because it shouldn't happen
			// every image/container should have this information
			log.L().Debug("this means either the hash or hash type is empty, returning an error because every image/container should have this information")
			return nil, createErr
		}

		images = append(images, newImage)
	}
	return images, nil
}
