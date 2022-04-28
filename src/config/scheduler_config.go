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
	RedisHostName      string
	RedisPort          string
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
		RedisHostName:      viper.GetString(cmd.RedisHostName),
		RedisPort:          viper.GetString(cmd.RedisPort),
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
