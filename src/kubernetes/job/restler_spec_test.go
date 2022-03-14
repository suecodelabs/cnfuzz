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
