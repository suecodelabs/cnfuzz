package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetImageName(t *testing.T) {
	imageName := "myregistry/someimage"
	testTags := []string{"latest", "debian-buster", "2.4.1", "slim"}
	testCase := fmt.Sprintf("%s:%s:%s:%s:%s", imageName, testTags[0], testTags[1], testTags[2], testTags[3])

	name, tags := GetImageName(testCase)
	assert.Equal(t, imageName, name)
	assert.Len(t, tags, len(tags))

}

func TestGetImageNameWithPort(t *testing.T) {
	// registry/image:tag1:tag2
	imageName := "localhost:1234/myImage"
	testTags := []string{"latest", "1.0"}
	testCase := fmt.Sprintf("%s:%s:%s", imageName, testTags[0], testTags[1])

	name, tags := GetImageName(testCase)
	assert.Equal(t, imageName, name)
	assert.Len(t, tags, len(testTags))
}

func TestSplitImageId(t *testing.T) {
	// docker-pullable://localhost:5000/imagename@sha256:5add8f7cf10b367af0fd4d9779a48891d9083ab56a691065421571b4d4cf4789

	testImageName := "registry.org:1234/myimage"
	testHashType := "sha256"
	testHash := "7471095e564d669dd964845e555b50752b30caa9d6d46e71d2e9278d63c57628"
	testCase := fmt.Sprintf("docker-pullable://%s@%s:%s", testImageName, testHashType, testHash)

	hash, hashType := SplitImageId(testCase)
	assert.Equal(t, testHashType, hashType)
	assert.Equal(t, testHash, hash)
}

func TestSplitImageId1(t *testing.T) {
	tests := []struct {
		name         string
		imageId      string
		wantHash     string
		wantHashType string
	}{
		{name: "empty-test", imageId: "", wantHash: "", wantHashType: ""},
		{name: "normal-test", imageId: "docker-pullable://localhost:5000/imagename@sha256:5add8f7cf10b367af0fd4d9779a48891d9083ab56a691065421571b4d4cf4789", wantHash: "5add8f7cf10b367af0fd4d9779a48891d9083ab56a691065421571b4d4cf4789", wantHashType: "sha256"},
		{name: "no-docker-prefix-test", imageId: "localhost:5000/imagename@sha256:abcdefghijklmnopqrstuvw", wantHash: "abcdefghijklmnopqrstuvw", wantHashType: "sha256"},
		{name: "common-registry-test", imageId: "docker-pullable://bla.com/imagename@sha256:abcdefghijklmnopqrstuvw", wantHash: "abcdefghijklmnopqrstuvw", wantHashType: "sha256"},
		{name: "short-hash-test", imageId: "docker-pullable://bla.com/imagename@sha256:a", wantHash: "a", wantHashType: "sha256"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHash, gotHashType := SplitImageId(tt.imageId)
			assert.Equalf(t, tt.wantHash, gotHash, "SplitImageId(%v)", tt.imageId)
			assert.Equalf(t, tt.wantHashType, gotHashType, "SplitImageId(%v)", tt.imageId)
		})
	}
}
