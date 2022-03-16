package config

import (
	"github.com/spf13/viper"
	"github.com/suecodelabs/cnfuzz/src/cmd"
	"github.com/suecodelabs/cnfuzz/src/log"
	apiv1 "k8s.io/api/core/v1"
)

// KubernetesFuzzConfig object that holds configuration for Kubernetes
type KubernetesFuzzConfig struct {
	TargetPodName      string
	TargetPodNamespace string
	ServiceAccountName string
	JobName            string
	RestlerJobName     string
	RestlerInitJobName string
	RestlerImage       string
	RestlerInitImage   string
	Namespace          string
	Image              string
}

// CreateKubernetesConfigWPod creates a KubernetesFuzzConfig from a pod object
func CreateKubernetesConfigWPod(pod *apiv1.Pod) *KubernetesFuzzConfig {
	namespace := getNamespace()
	return &KubernetesFuzzConfig{
		TargetPodName:      pod.Name,
		TargetPodNamespace: pod.Namespace,
		ServiceAccountName: "cnfuzz-job",
		JobName:            "cnfuzz-" + pod.Name,
		RestlerJobName:     "cnfuzz-restler-" + pod.Name,
		RestlerImage:       viper.GetString(cmd.RestlerImageFlag),
		Namespace:          namespace,
		Image:              viper.GetString(cmd.JobImageFlag),
	}
}

// CreateKubernetesConfig creates a KubernetesFuzzConfig for a pod with its name and namespace
func CreateKubernetesConfig(targetPodName string, namespace string) *KubernetesFuzzConfig {
	return &KubernetesFuzzConfig{
		TargetPodName:      targetPodName,
		ServiceAccountName: "cnfuzz-job",
		JobName:            "cnfuzz-" + targetPodName,
		RestlerInitJobName: "cnfuzz-restler-init-" + targetPodName,
		RestlerJobName:     "cnfuzz-restler-" + targetPodName,
		RestlerInitImage:   viper.GetString(cmd.RestlerInitImageFlag),
		RestlerImage:       viper.GetString(cmd.RestlerImageFlag),
		Namespace:          namespace,
		Image:              viper.GetString(cmd.JobImageFlag),
	}
}

// getNamespace function that gets the home namespace from viper and checks if it's valid
func getNamespace() string {
	namespace := viper.GetString(cmd.HomeNamespaceFlag)
	if len(namespace) <= 0 {
		log.L().Fatalf("\"%s\" is not a valid namespace", namespace)
	}
	return namespace
}
