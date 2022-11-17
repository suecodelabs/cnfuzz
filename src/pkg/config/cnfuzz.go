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
	"regexp"
)

const imageRegex = "[a-z0-9]([-a-z0-9]*[a-z0-9])?"

type CnFuzzConfig struct {
	Namespace            string                `yaml:"namespace" `
	OnlyFuzzMarked       bool                  `yaml:"only_fuzz_marked"`
	CacheSolution        string                `yaml:"cache_solution"`
	ConfigmapName        string                `yaml:"configmap_name"`
	RestlerWrapperConfig *RestlerWrapperConfig `yaml:"restlerwrapper"`
	RedisConfig          *RedisConfig          `yaml:"redis"`
	AuthConfig           *AuthConfig           `yaml:"auth"`
	S3Config             *S3Config             `yaml:"s3"`
}

type ImageConfig struct {
	Image      string `yaml:"image"`
	PullPolicy string `yaml:"pullPolicy"`
	Tag        string `yaml:"tag"`
}

func (cnf ImageConfig) GetImage() string {
	var containerImage string
	if len(cnf.Tag) > 0 {
		containerImage = fmt.Sprintf("%s:%s", cnf.Image, cnf.Tag)
	} else {
		containerImage = cnf.Image
	}
	return containerImage
}

func (cnf ImageConfig) Validate() (bool, error) {
	if len(cnf.Image) == 0 {
		return false, fmt.Errorf("image is empty")
	}
	return regexp.MatchString(imageRegex, cnf.GetImage())
}

type RestlerWrapperConfig struct {
	ImageConfig    ImageConfig    `yaml:"image"`
	RestlerConfig  *RestlerConfig `yaml:"restler"`
	ServiceAccount string         `yaml:"service_account"`
}

type RestlerConfig struct {
	TimeBudget      string `yaml:"time_budget"`
	CpuLimit        string `yaml:"cpu_limit"`
	MemoryLimit     string `yaml:"memory_limit"`
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

	if config.RestlerWrapperConfig == nil {
		return nil, fmt.Errorf("give config file doesn't contain configuration for the restlerwrapper")
	}
	v, err := config.RestlerWrapperConfig.ImageConfig.Validate()
	if !v || err != nil {
		if err != nil {
			l.Error(err, "given restler wrapper image is invalid")
			return nil, err
		}
		return nil, fmt.Errorf("given restler wrapper image is invalid, needs to match '%s'", imageRegex)
	}

	return config, nil
}
