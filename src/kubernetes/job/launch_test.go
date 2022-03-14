package job

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/suecodelabs/cnfuzz/src/config"
	"github.com/suecodelabs/cnfuzz/src/discovery"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"net/url"
	"testing"
)

func TestLaunchK8sJob(t *testing.T) {
	clientSet := fake.NewSimpleClientset()
	testNamespaceName := "my-test-ns"
	testJobName := "my-job"
	config := &config.KubernetesFuzzConfig{
		Namespace: testNamespaceName,
		JobName:   testJobName,
	}

	LaunchFuzzJob(clientSet, config)
	result, err := clientSet.BatchV1().Jobs(testNamespaceName).Get(context.TODO(), testJobName, metav1.GetOptions{})
	if assert.NoError(t, err) {
		assert.Equal(t, testJobName, result.Name)
	}
}

func TestLaunchRestlerJob(t *testing.T) {
	clientSet := fake.NewSimpleClientset()
	testNamespaceName := "my-test-ns"
	testJobName := "my-restler-job"
	testUri, _ := url.Parse("http://testservice:8080")
	fuzzConf := &config.FuzzConfig{
		ApiDescription: &discovery.WebApiDescription{
			BaseUrl: *testUri,
		},
		KubernetesConfig: &config.KubernetesFuzzConfig{
			RestlerImage:   "mcr.microsoft.com/restlerfuzzer/restler:v7.4.0",
			RestlerJobName: testJobName,
			Namespace:      testNamespaceName,
		},
	}

	_, launchErr := LaunchRestlerJob(clientSet, fuzzConf, &v1.Pod{ObjectMeta: metav1.ObjectMeta{
		Name:      "target-api",
		Namespace: testNamespaceName,
	}})
	assert.NoError(t, launchErr)
	result, err := clientSet.BatchV1().Jobs(testNamespaceName).Get(context.TODO(), testJobName, metav1.GetOptions{})
	if assert.NoError(t, err) {
		assert.Equal(t, testJobName, result.Name)
	}
}
