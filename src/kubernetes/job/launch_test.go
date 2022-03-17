package job

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/suecodelabs/cnfuzz/src/auth"
	"github.com/suecodelabs/cnfuzz/src/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestLaunchK8sJob(t *testing.T) {
	clientSet := fake.NewSimpleClientset()
	testNamespaceName := "my-test-ns"
	testJobName := "my-job"
	config := &config.SchedulerConfig{
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

	restlerConf := &config.FuzzerConfig{
		JobName:   testJobName,
		Namespace: testNamespaceName,
		Target: config.FuzzerTarget{
			IP:   "10.0.0.1",
			Port: "8080",
		},
	}

	testTokenSource := testTSource{}

	_, launchErr := LaunchRestlerJob(clientSet, restlerConf, testTokenSource)
	assert.NoError(t, launchErr)
	result, err := clientSet.BatchV1().Jobs(testNamespaceName).Get(context.TODO(), testJobName, metav1.GetOptions{})
	if assert.NoError(t, err) {
		assert.Equal(t, testJobName, result.Name)
	}
}

type testTSource struct {
}

func (testTSource) Token() (*auth.Token, error) {
	return &auth.Token{}, nil
}
