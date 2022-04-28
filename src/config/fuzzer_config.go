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
	"github.com/suecodelabs/cnfuzz/src/discovery"
	v1 "k8s.io/api/core/v1"
)

// FuzzerConfig config that is used to configure the RESTler fuzzer
type FuzzerConfig struct {
	JobName                string
	InitJobName            string
	Namespace              string
	Image                  string
	InitImage              string
	TimeBudget             string
	RestlerTelemetryOptOut string
	CpuLimit               int64
	MemoryLimit            int64
	CpuRequest             int64
	MemoryRequest          int64
	DiscoveryDocLocation   string
	Target                 FuzzerTarget
	S3Config               S3Config
}

// FuzzerTarget configuration for the fuzzing target
type FuzzerTarget struct {
	PodName   string
	Namespace string
	IP        string
	Port      string
	Scheme    string // http, https
}

// NewFuzzerConfig constructor for FuzzerConfig
func NewFuzzerConfig(apiDesc *discovery.WebApiDescription, targetPod *v1.Pod) *FuzzerConfig {
	ns := viper.GetString(cmd.HomeNamespaceFlag)
	return &FuzzerConfig{
		JobName:              "cnfuzz-restler-" + targetPod.Name,
		InitJobName:          "cnfuzz-restler-init-" + targetPod.Name,
		Image:                viper.GetString(cmd.RestlerImageFlag),
		Namespace:            ns,
		InitImage:            viper.GetString(cmd.RestlerInitImageFlag),
		TimeBudget:           viper.GetString(cmd.RestlerTimeBudget),
		CpuLimit:             viper.GetInt64(cmd.RestlerCpuLimit),
		MemoryLimit:          viper.GetInt64(cmd.RestlerMemoryLimit),
		CpuRequest:           viper.GetInt64(cmd.RestlerCpuRequest),
		MemoryRequest:        viper.GetInt64(cmd.RestlerMemoryRequest),
		DiscoveryDocLocation: apiDesc.DiscoveryDoc.String(),
		Target: FuzzerTarget{
			PodName:   targetPod.Name,
			Namespace: targetPod.Namespace,
			IP:        targetPod.Status.PodIP,
			Port:      apiDesc.BaseUrl.Port(),
			Scheme:    apiDesc.BaseUrl.Scheme,
		},
		S3Config: *CreateS3Config("cnfuzz-aws-cli-"+targetPod.Name, ns),
	}
}
