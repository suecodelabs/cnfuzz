package model

type ImageFuzzStatus int

const (
	Unknown ImageFuzzStatus = iota
	NotFuzzed
	Fuzzed
)

// ContainerImageVersion a container image can have multiple versions
// different version can be a different set of tags or just a different hash
// only the hash should be used to differentiate between images, because the only trustworthy property
type ContainerImageVersion struct {
	Hash     string
	HashType string
	Tags     []string
	Status   ImageFuzzStatus
}

// CreateContainerImageVersion constructor for ContainerImageVersion
func CreateContainerImageVersion(hash string, hashType string, tags []string, status ImageFuzzStatus) *ContainerImageVersion {
	return &ContainerImageVersion{
		Hash:     hash,
		HashType: hashType,
		Tags:     tags,
		Status:   status,
	}
}
