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

const ConfigMapName = "cnfuzz_config"

// SchedulerConfig object that holds configuration for Kubernetes
type SchedulerConfig struct {
	ServiceAccountName string `yaml:"service_account_name"`
	RedisHostName      string `yaml:"redis_host_name"`
	RedisPort          string `yaml:"redis_port"`
	Namespace          string `yaml:"namespace"`
	OnlyFuzzMarked     bool   `yaml:"only_fuzz_marked"`
}

type RestlerConfig struct {
	InitImage     string `yaml:"init_image"`
	Image         string `yaml:"image"`
	TimeBudget    string `yaml:"time_budget"`
	CpuLimit      string `yaml:"cpu_limit"`
	MemoryLimit   string `yaml:"memory_limit"`
	CpuRequest    string `yaml:"cpu_request"`
	MemoryRequest string `yaml:"memory_request"`
	CacheSolution string `yaml:"cache_solution"`
}

type AuthConfig struct {
	Username string `yaml:"username"`
	Secret   string `yaml:"secret"`
}

type S3Config struct {
	ContainerName string `yaml:"container_name"`
	Namespace     string `yaml:"namespace"`
	EndpointUrl   string `yaml:"endpoint_url"`
	ReportBucket  string `yaml:"report_bucket"`
	Image         string `yaml:"image"`
	AccessKey     string `yaml:"access_key"`
	SecretKey     string `yaml:"secret_key"`
}

// CRestlerCpuLimit,     reateSchedulerConfigWPod creates a SchedulerConfig from a pod object
/* fRestlerMemoryLimit,  unc CreateSchedulerConfigWPod(pod *apiv1.Pod) *SchedulerConfig {
	RestlerCpuRequest,   namespace := getNamespace()
	RestlerMemoryRequest,return &SchedulerConfig{
		TargetPodName:      pod.Name,
		TargetPodNamespace: pod.Namespace,
		ServiceAccountName: "cnfuzz-scheduler",
		JobName:            "cnfuzz-" + pod.Name,
		RedisHostName:      viper.GetString(flags.RedisHostName),
		RedisPort:          viper.GetString(flags.RedisPort),
		Namespace:          namespace,
		Image:              viper.GetString(flags.SchedulerImageFlag),
	}
} */

// getNamespace function that gets the home namespace from viper and checks if it's valid
/* func getNamespace() string {
	namespace := viper.GetString(flags.HomeNamespaceFlag)
	if len(namespace) <= 0 {
		log.L().Fatalf("\"%s\" is not a valid namespace", namespace)
	}
	return namespace
} */
