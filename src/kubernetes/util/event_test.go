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
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestIsKubernetesEvent(t *testing.T) {
	nameTestCases := []string{"some-object", "kubernetes", "some-object"}
	namespaceTestCases := []string{"kube-system", "default", "default"}
	results := []bool{true, true, false}
	for i := 0; i < len(nameTestCases); i++ {
		testMeta := &v1.ObjectMeta{
			Name:      nameTestCases[i],
			Namespace: namespaceTestCases[i],
		}
		result := IsKubernetesEvent(testMeta)
		assert.Equal(t, results[i], result)
	}
}

func TestIsFuzzerEvent(t *testing.T) {
	testCases := []string{"cnfuzz-myapiimage-123864", "myapiimage-123864"}
	results := []bool{true, false}

	for i := 0; i < len(testCases); i++ {
		testMeta := &v1.ObjectMeta{
			Name: testCases[i],
		}
		result := IsFuzzerEvent(testMeta)
		assert.Equal(t, results[i], result)
	}
}
