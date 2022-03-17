package model

import (
	kutil "github.com/suecodelabs/cnfuzz/src/kubernetes/util"
	"github.com/suecodelabs/cnfuzz/src/log"
	apiv1 "k8s.io/api/core/v1"
)

// ContainerImage object that holds information about a container image
type ContainerImage struct {
	Id       string
	Name     string
	Versions []ContainerImageVersion
}

// CreateContainerImagesFromPod extract all container images from a pod
func CreateContainerImagesFromPod(pod *apiv1.Pod) ([]ContainerImage, error) {
	logger := log.L()
	var images []ContainerImage
mainloop:
	for _, status := range pod.Status.ContainerStatuses {
		if len(status.ImageID) == 0 || len(status.Image) == 0 {
			// TODO are there other image names or ID's that we can't parse/check?
			logger.Warnf("image ID \"%s\" and image name \"%s\" are invalid and can't be checked", status.ImageID, status.Image)
			continue
		}

		hash, hashType := kutil.SplitImageId(status.ImageID)
		name, tags := kutil.GetImageName(status.Image)

		// Look for duplicate image hashes/versions
		for _, image := range images {
			if image.Name == name {
				// Found an already existing image
				// Add possible different versions
				for _, version := range image.Versions {
					if version.Hash != hash {
						newVersion := ContainerImageVersion{
							Hash:     hash,
							HashType: hashType,
							Status:   Unknown,
						}
						image.Versions = append(image.Versions, newVersion)
					}
				}
				// Image already exists in the image array, so avoid creating a duplicate
				continue mainloop
			}
		}

		var imageVersions = make([]ContainerImageVersion, 1)
		// Container can only have one image
		if len(imageVersions) > 1 {
			logger.Infof("container %s contains more then one image(?), fuzzer will only remember the first one", status.Name)
		}
		imageVersions[0] = *CreateContainerImageVersion(hash, hashType, tags, Unknown) // Unknown is the fuzz status (has this image already been fuzzed)
		newImage := ContainerImage{
			Id:       status.ImageID,
			Name:     name,
			Versions: imageVersions,
		}

		images = append(images, newImage)
	}
	return images, nil
}
