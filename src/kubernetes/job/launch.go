package job

import (
	"context"
	"github.com/suecodelabs/cnfuzz/src/auth"

	"github.com/suecodelabs/cnfuzz/src/config"
	"github.com/suecodelabs/cnfuzz/src/log"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// LaunchFuzzJob launches a job that retrieves the OpenAPI doc and kicks of a RESTler job
func LaunchFuzzJob(clientset kubernetes.Interface, kubeConfig *config.SchedulerConfig) (createdJob *batchv1.Job, err error) {
	jobs := clientset.BatchV1().Jobs(kubeConfig.Namespace)

	jobSpec := createSchedulerJob(kubeConfig)

	result, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	return result, err
}

// LaunchRestlerJob starts a job with the RESTler fuzzer
func LaunchRestlerJob(clientset kubernetes.Interface, restlerConfig *config.FuzzerConfig, tokenSource auth.ITokenSource) (restlerJob *batchv1.Job, err error) {
	// Start the RESTler container as a job
	restlerJobSpec := createRestlerJob(restlerConfig, tokenSource)
	restlerJob, jErr := clientset.BatchV1().Jobs(restlerConfig.Namespace).Create(context.TODO(), restlerJobSpec, metav1.CreateOptions{})
	if jErr != nil {
		log.L().Errorf("failed to create restler job: %+v", jErr)
		return nil, jErr
	}
	return restlerJobSpec, nil
}
