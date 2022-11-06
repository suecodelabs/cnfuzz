/*
 * Copyright 2022 Sue B.V.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"fmt"
	"github.com/suecodelabs/cnfuzz/src/pkg/logger"
	"gopkg.in/yaml.v2"
)

type CnFuzzConfig struct {
	Namespace      string          `yaml:"namespace" `
	OnlyFuzzMarked bool            `yaml:"only_fuzz_marked"`
	CacheSolution  string          `yaml:"cache_solution"`
	ConfigmapName  string          `yaml:"configmap_name"`
	RestlerWrapper *RestlerWrapper `yaml:"restlerwrapper"`
	RedisConfig    *RedisConfig    `yaml:"redis"`
	AuthConfig     *AuthConfig     `yaml:"auth"`
	S3Config       *S3Config       `yaml:"s3"`
}

type ImageConfig struct {
	Image      string `yaml:"image"`
	PullPolicy string `yaml:"pullPolicy"`
	Tag        string `yaml:"tag"`
}

type RestlerWrapper struct {
	ImageConfig    `yaml:"image"`
	*RestlerConfig `yaml:"restler"`
}

type RestlerConfig struct {
	TimeBudget      string `yaml:"time_budget"`
	CpuLimit        int64  `yaml:"cpu_limit"`
	MemoryLimit     int64  `yaml:"memory_limit"`
	CpuRequest      string `yaml:"cpu_request"`
	MemoryRequest   string `yaml:"memory_request"`
	TelemetryOptOut string `yaml:"telemetry_opt_out"`
}

type RedisConfig struct {
	HostName string `yaml:"host_name"`
	Port     string `yaml:"port"`
}

type AuthConfig struct {
	Username string `yaml:"username"`
	Secret   string `yaml:"secret"`
}

type S3Config struct {
	EndpointUrl  string `yaml:"endpoint_url"`
	ReportBucket string `yaml:"report_bucket"`
	AccessKey    string `yaml:"access_key"`
	SecretKey    string `yaml:"secret_key"`
}

func LoadCnFuzzConfig(l logger.Logger, configFile string, printFile bool) (*CnFuzzConfig, error) {
	config := &CnFuzzConfig{}

	data, err := loadFile(l, configFile, printFile)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(*data, config); err != nil {
		return nil, err
	}

	if config.RestlerWrapper == nil {
		return nil, fmt.Errorf("give config file doesn't contain configuration for the restlerwrapper")
	}

	return config, nil
}
