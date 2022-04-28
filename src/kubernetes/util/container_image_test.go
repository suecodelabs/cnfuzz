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
		{name: "just-hash-and-type-test", imageId: "sha256:46aa7ac68facde8183f9df7c059b7e7e1ac45ae73157512d9a1ce57ae4fe5eec", wantHash: "46aa7ac68facde8183f9df7c059b7e7e1ac45ae73157512d9a1ce57ae4fe5eec", wantHashType: "sha256"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHash, gotHashType := SplitImageId(tt.imageId)
			assert.Equalf(t, tt.wantHash, gotHash, "SplitImageId(%v)", tt.imageId)
			assert.Equalf(t, tt.wantHashType, gotHashType, "SplitImageId(%v)", tt.imageId)
		})
	}
}
