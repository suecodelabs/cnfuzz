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

package job

// TODO this test breaks the Golang testing package???
/* func TestCreateRestlerJob(t *testing.T) {
	testUri, _ := url.Parse("http://testservice:8080")
	testKubernetesConfig := &config.KubernetesFuzzConfig{
		RestlerJobName: "my-test-restler-job",
	}
	testConfig := &config.FuzzConfig{
		// ApiDescription.BaseUrl
		ApiDescription: &discovery.WebApiDescription{
			BaseUrl: *testUri,
		},

		KubernetesConfig: testKubernetesConfig,
	}

	testTargetPod := &v1.Pod{}

	// Start the RESTler container as a job
	restlerJobSpec := createRestlerJob(testConfig, testTargetPod)
	assert.Equal(t, restlerJobSpec.Annotations["cnfuzz/ignore"], "true")
	assert.Equal(t, restlerJobSpec.Name, testKubernetesConfig.RestlerJobName)
	assert.Equal(t, restlerJobSpec.Namespace, testKubernetesConfig.Namespace)
	firstVol := restlerJobSpec.Spec.Template.Spec.Volumes[0]
	firstContainerSpec := restlerJobSpec.Spec.Template.Spec.Containers[0]
	if assert.NotNil(t, firstVol) {
		firstVolName := firstVol.Name
		if assert.NotNil(t, firstContainerSpec) {
			firstVolMount := firstContainerSpec.VolumeMounts[0]
			if assert.NotNil(t, firstVolMount) {
				assert.Equal(t, firstVolName, firstVolMount.Name)
			}
		}
	}
} */
