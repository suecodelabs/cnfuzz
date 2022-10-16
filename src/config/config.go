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
	"github.com/suecodelabs/cnfuzz/src/logger"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

// CnFuzzConfig object that holds configuration for Kubernetes
type CnFuzzConfig struct {
	Namespace      string `yaml:"namespace" `
	OnlyFuzzMarked bool   `yaml:"only_fuzz_marked"`
	CacheSolution  string `yaml:"cache_solution"`
	*RedisConfig   `yaml:"redis"`
	*RestlerConfig `yaml:"restler"`
	*AuthConfig    `yaml:"auth"`
	*S3Config      `yaml:"s3"`
}

type RedisConfig struct {
	HostName string `yaml:"host_name"`
	Port     string `yaml:"port"`
}

type RestlerConfig struct {
	InitImage       string `yaml:"init_image"`
	Image           string `yaml:"image"`
	TimeBudget      string `yaml:"time_budget"`
	CpuLimit        int64  `yaml:"cpu_limit"`
	MemoryLimit     int64  `yaml:"memory_limit"`
	CpuRequest      string `yaml:"cpu_request"`
	MemoryRequest   string `yaml:"memory_request"`
	TelemetryOptOut string `yaml:"telemetry_opt_out"`
}

type AuthConfig struct {
	Username string `yaml:"username"`
	Secret   string `yaml:"secret"`
}

type S3Config struct {
	EndpointUrl  string          `yaml:"endpoint_url"`
	ReportBucket string          `yaml:"report_bucket"`
	AccessKey    string          `yaml:"access_key"`
	SecretKey    string          `yaml:"secret_key"`
	Sidecar      S3SidecarConfig `yaml:"sidecar"`
}

type S3SidecarConfig struct {
	Image         string `yaml:"image"`
	ContainerName string `yaml:"container_name"`
}

func LoadConfigFile(l logger.Logger, configFile string, printFile bool) *CnFuzzConfig {
	if configFile == "" {
		configFile = filepath.Join()
		if !pathExists(configFile) {
			// return nil
		}
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		l.V(logger.ImportantLevel).Error(err, "failed to read config file", "configFile", configFile)
		os.Exit(1)
	}

	if printFile {
		l.V(logger.DebugLevel).Info("\n" + string(data[:]))
	}

	config := &CnFuzzConfig{}
	if err := yaml.Unmarshal(data, config); err != nil {
		l.V(logger.ImportantLevel).Error(err, "given config file is invalid")
		os.Exit(1)
	}

	if config.RedisConfig == nil {

	}

	if config.RestlerConfig == nil {

	}

	if config.AuthConfig == nil {

	}

	if config.S3Config == nil {

	}

	return config
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
