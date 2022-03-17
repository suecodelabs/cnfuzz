package config

import (
	"github.com/spf13/viper"
	"github.com/suecodelabs/cnfuzz/src/cmd"
	"github.com/suecodelabs/cnfuzz/src/log"
	apiv1 "k8s.io/api/core/v1"
)

// SchedulerConfig object that holds configuration for Kubernetes
type SchedulerConfig struct {
	TargetPodName      string
	TargetPodNamespace string
	ServiceAccountName string
	JobName            string
	Namespace          string
	Image              string
}

// CreateSchedulerConfigWPod creates a SchedulerConfig from a pod object
func CreateSchedulerConfigWPod(pod *apiv1.Pod) *SchedulerConfig {
	namespace := getNamespace()
	return &SchedulerConfig{
		TargetPodName:      pod.Name,
		TargetPodNamespace: pod.Namespace,
		ServiceAccountName: "cnfuzz-scheduler",
		JobName:            "cnfuzz-" + pod.Name,
		Namespace:          namespace,
		Image:              viper.GetString(cmd.SchedulerImageFlag),
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
