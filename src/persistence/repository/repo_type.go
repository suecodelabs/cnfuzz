package repository

import "fmt"

type RepoType int

const (
	Redis RepoType = iota
	InMemory
)

var RepoTypes = [2]string{"redis", "in_memory"}

func (s RepoType) String() string {
	return RepoTypes[s]
}

func RepoTypeFromString(value string) (RepoType, error) {
	for i := 0; i < len(RepoTypes); i++ {
		if RepoTypes[i] == value {
			return RepoType(i), nil
		}
	}
	return RepoType(0), fmt.Errorf("unkown repo type")
}
