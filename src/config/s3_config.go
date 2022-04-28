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
)

type S3Config struct {
	ContainerName string
	Namespace     string
	EndpointUrl   string
	ReportBucket  string
	Image         string
	AccessKey     string
	SecretKey     string
}

func CreateS3Config(jobName string, namespace string) *S3Config {
	s3Config := &S3Config{
		ContainerName: jobName,
		Namespace:     namespace,
		EndpointUrl:   viper.GetString(cmd.S3EndpointUrlFlag),
		ReportBucket:  viper.GetString(cmd.S3ReportBucket),
		Image:         "amazon/aws-cli",

		AccessKey: viper.GetString(cmd.S3AccessKey),
		SecretKey: viper.GetString(cmd.S3SecretKey),
	}
	if len(s3Config.ReportBucket) == 0 {
		log.L().Fatal("no S3 bucket given to store reports in!")
	}
	return s3Config
}
