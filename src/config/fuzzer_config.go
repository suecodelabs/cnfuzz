package config

import (
	"github.com/spf13/viper"
	"github.com/suecodelabs/cnfuzz/src/cmd"
	"github.com/suecodelabs/cnfuzz/src/discovery"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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
	CpuLimits      resource.Quantity
	CpuRequests    resource.Quantity
	MemoryLimits   resource.Quantity
	MemoryRequests resource.Quantity
}

// NewFuzzerConfig constructor for FuzzerConfig
func NewFuzzerConfig(apiDesc *discovery.WebApiDescription, targetPod *v1.Pod) *FuzzerConfig {
	res := FuzzerJobResources{
		CpuLimits:      resource.MustParse(viper.GetString(cmd.RestlerCpuLimits)),
		CpuRequests:    resource.MustParse(viper.GetString(cmd.RestlerCpuRequests)),
		MemoryLimits:   resource.MustParse(viper.GetString(cmd.RestlerMemoryLimits)),
		MemoryRequests: resource.MustParse(viper.GetString(cmd.RestlerMemoryRequests)),
	}

	return &FuzzerConfig{
		JobName:              "cnfuzz-restler-" + targetPod.Name,
		InitJobName:          "cnfuzz-restler-init-" + targetPod.Name,
		Image:                viper.GetString(cmd.RestlerImageFlag),
		Namespace:            viper.GetString(cmd.HomeNamespaceFlag),
		InitImage:            viper.GetString(cmd.RestlerInitImageFlag),
		TimeBudget:           viper.GetString(cmd.RestlerTimeBudget),
		DiscoveryDocLocation: apiDesc.DiscoveryDoc.String(),
		Resources:            res,
		Target: FuzzerTarget{
			PodName:   targetPod.Name,
			Namespace: targetPod.Namespace,
			IP:        targetPod.Status.PodIP,
			Port:      apiDesc.BaseUrl.Port(),
			Scheme:    apiDesc.BaseUrl.Scheme,
		},
	}
}
