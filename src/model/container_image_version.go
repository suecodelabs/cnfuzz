package model

type ImageFuzzStatus int

const (
	Unknown ImageFuzzStatus = iota
	NotFuzzed
	Fuzzed
)

type ContainerImageVersion struct {
	Hash     string
	HashType string
	Tags     []string
	Status   ImageFuzzStatus
}

func CreateContainerImageVersion(hash string, hashType string, tags []string, status ImageFuzzStatus) *ContainerImageVersion {
	return &ContainerImageVersion{
		Hash:     hash,
		HashType: hashType,
		Tags:     tags,
		Status:   status,
	}
}
