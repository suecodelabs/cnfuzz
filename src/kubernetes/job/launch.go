package job

import (
	"context"

	"github.com/suecodelabs/cnfuzz/src/config"
	"github.com/suecodelabs/cnfuzz/src/log"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func LaunchFuzzJob(clientset kubernetes.Interface, kubeConfig *config.KubernetesFuzzConfig) (createdJob *batchv1.Job, err error) {
	jobs := clientset.BatchV1().Jobs(kubeConfig.Namespace)

	jobSpec := GetJobSpec(kubeConfig)

	result, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	return result, err
}

func LaunchRestlerJob(clientset kubernetes.Interface, fuzzConfig *config.FuzzConfig, targetPod *v1.Pod) (restlerJob *batchv1.Job, err error) {
	// Start the RESTler container as a job
	restlerJobSpec := createRestlerJob(fuzzConfig, targetPod)
	restlerJob, jErr := clientset.BatchV1().Jobs(fuzzConfig.KubernetesConfig.Namespace).Create(context.TODO(), restlerJobSpec, metav1.CreateOptions{})
	if jErr != nil {
		log.L().Errorf("failed to create restler job: %+v", jErr)
		return nil, jErr
	}
	return restlerJobSpec, nil
}
