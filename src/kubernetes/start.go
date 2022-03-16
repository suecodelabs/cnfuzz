package kubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/suecodelabs/cnfuzz/src/cmd"
	"github.com/suecodelabs/cnfuzz/src/config"
	"github.com/suecodelabs/cnfuzz/src/discovery/openapi"
	"github.com/suecodelabs/cnfuzz/src/kubernetes/job"
	"github.com/suecodelabs/cnfuzz/src/kubernetes/util"
	"github.com/suecodelabs/cnfuzz/src/log"
	"github.com/suecodelabs/cnfuzz/src/persistence/repository"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// FuzzPod fuzz a pod from a Pod object
func FuzzPod(pod *v1.Pod) error {
	insideCluster := viper.GetBool(cmd.InsideClusterFlag)
	clientset, err := CreateClientSet(insideCluster)
	if err != nil {
		return fmt.Errorf("error while getting kubernetes clientset: %w", err)
	}

	return fuzzPod(clientset, pod)
}

// FuzzPodWithName fuzz a pod from just its name and namespace
func FuzzPodWithName(namespace string, podName string) (err error) {
	insideCluster := viper.GetBool(cmd.InsideClusterFlag)
	clientset, err := CreateClientSet(insideCluster)
	if err != nil {
		return fmt.Errorf("error while getting kubernetes clientset: %w", err)
	}

	if !util.NamespaceExist(clientset, namespace) {
		return fmt.Errorf("namespace %s doesn't exist or is invalid", namespace)
	}

	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to find target pod %s inside namespace %s: %w", podName, namespace, err)
	}

	return fuzzPod(clientset, pod)
}

// StartInformers start informers that listen for Kubernetes events and let the EventHandler react on the events
func StartInformers(repos *repository.Repositories) (err error) {
	insideCluster := viper.GetBool(cmd.InsideClusterFlag)
	clientset, err := CreateClientSet(insideCluster)
	if err != nil {
		return fmt.Errorf("error while getting kubernetes clientset: %w", err)
	}

	myEventHandler := NewEventHandler(clientset, repos)

	factory := informers.NewSharedInformerFactory(clientset, time.Hour*24)
	podInformer := factory.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(myEventHandler)

	log.L().Info("start listening for events ...")

	stopChan := make(chan struct{})
	defer close(stopChan)
	factory.Start(stopChan)
	if !cache.WaitForCacheSync(stopChan, podInformer.HasSynced) {
		log.L().Error("waiting for cache to populate from Kubernetes was unsuccessful")
		return
	}
	select {}
}

// fuzzPod start the fuzzing process for a pod
func fuzzPod(clientSet kubernetes.Interface, pod *v1.Pod) error {
	annos := GetAnnotations(&pod.ObjectMeta)
	annos.SetConfigRegister()

	ip := pod.Status.PodIP
	var ports []int32
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			ports = append(ports, port.ContainerPort)
		}
	}

	var oaLocs []string
	if len(annos.OpenApiDocLocation) > 0 {
		oaLocs = append(oaLocs, annos.OpenApiDocLocation)
	} else {
		oaLocs = openapi.GetCommonOpenApiLocations()
	}

	apiDesc, err := openapi.TryGetOpenApiDoc(ip, ports, oaLocs)
	if err != nil {
		return fmt.Errorf("error while retrieving OpenAPI document from target %s: %w", pod.Name, err)
	}

	fuzzConfig := config.NewFuzzConfig(apiDesc, pod.Name, pod.Namespace)

	restlerJob, restlerErr := job.LaunchRestlerJob(clientSet, fuzzConfig, pod)
	if restlerErr != nil {
		return fmt.Errorf("error while starting restler pod: %v", restlerErr)
	}

	// TODO start monitoring pod

	log.L().Debug("waiting for RESTler job to complete")
	// Long timeout because restler jobs can take a long time
	waitErr := util.WaitForJobReady(clientSet, restlerJob.Name, restlerJob.Namespace, time.Hour*2)
	if waitErr != nil {
		// We dont want to leave the config map hanging around, so remove it
		return fmt.Errorf("error while waiting for job to finish: %v", waitErr)
	}
	// TODO tell monitor to stop and get results

	log.L().Debug("job complete!")
	return nil
}
