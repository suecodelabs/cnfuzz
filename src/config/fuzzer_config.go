package config

import (
	"github.com/spf13/viper"
	"github.com/suecodelabs/cnfuzz/src/cmd"
	"github.com/suecodelabs/cnfuzz/src/discovery"
	v1 "k8s.io/api/core/v1"
)

// FuzzerConfig config that is used to configure the RESTler fuzzer
type FuzzerConfig struct {
	JobName              string
	InitJobName          string
	Namespace            string
	Image                string
	InitImage            string
	TimeBudget           string
	DiscoveryDocLocation string
	Target               FuzzerTarget
	Resources            FuzzerJobResources
}

// FuzzerTarget configuration for the fuzzing target
type FuzzerTarget struct {
	PodName   string
	Namespace string
	IP        string
	Port      string
	Scheme    string // http, https
}

// FuzzerJobResources contains resource limits/requests for all fuzzer jobs
type FuzzerJobResources struct {
	CpuLimits      string
	CpuRequests    string
	MemoryLimits   string
	MemoryRequests string
}

// NewFuzzerConfig constructor for FuzzerConfig
func NewFuzzerConfig(apiDesc *discovery.WebApiDescription, targetPod *v1.Pod) *FuzzerConfig {
	return &FuzzerConfig{
		JobName:              "cnfuzz-restler-" + targetPod.Name,
		InitJobName:          "cnfuzz-restler-init-" + targetPod.Name,
		Image:                viper.GetString(cmd.RestlerImageFlag),
		Namespace:            viper.GetString(cmd.HomeNamespaceFlag),
		InitImage:            viper.GetString(cmd.RestlerInitImageFlag),
		TimeBudget:           viper.GetString(cmd.RestlerTimeBudget),
		DiscoveryDocLocation: apiDesc.DiscoveryDoc.String(),
		Resources: FuzzerJobResources{
			CpuLimits:      viper.GetString(cmd.RestlerCpuLimits),
			CpuRequests:    viper.GetString(cmd.RestlerCpuRequests),
			MemoryLimits:   viper.GetString(cmd.RestlerMemoryLimits),
			MemoryRequests: viper.GetString(cmd.RestlerMemoryRequests),
		},
		Target: FuzzerTarget{
			PodName:   targetPod.Name,
			Namespace: targetPod.Namespace,
			IP:        targetPod.Status.PodIP,
			Port:      apiDesc.BaseUrl.Port(),
			Scheme:    apiDesc.BaseUrl.Scheme,
		},
	}
}
