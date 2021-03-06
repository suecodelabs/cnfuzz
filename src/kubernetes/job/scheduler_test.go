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

package job

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/suecodelabs/cnfuzz/src/cmd"
	"github.com/suecodelabs/cnfuzz/src/config"
	"testing"
)

func TestGetJobSpec(t *testing.T) {
	testConfig := &config.SchedulerConfig{
		JobName:   "test-job",
		Image:     "test-img",
		Namespace: "test",
	}
	result := createSchedulerJob(testConfig)
	assert.Equal(t, testConfig.JobName, result.Name)
	assert.Equal(t, testConfig.Namespace, result.Namespace)
	assert.Equal(t, "true", result.Annotations["cnfuzz/ignore"])
	assert.Equal(t, "cnfuzz-job", result.Spec.Template.Spec.ServiceAccountName)
	firstContainer := result.Spec.Template.Spec.Containers[0]
	if assert.NotNil(t, firstContainer) {
		assert.Equal(t, testConfig.JobName, firstContainer.Name)
		assert.Equal(t, testConfig.Image, firstContainer.Image)
	}
}

func TestBuildJobArgs(t *testing.T) {
	testConfig := &config.SchedulerConfig{
		TargetPodName:      "test-target",
		TargetPodNamespace: "test-target-namespace",
	}
	viper.Set(cmd.IsDebug, true)
	username := "user123"
	secret := "@Welcome123"
	targetNamespace := "the-test-namespace"
	viper.Set(cmd.AuthUsername, username)
	viper.Set(cmd.AuthSecretFlag, secret)
	viper.Set(cmd.HomeNamespaceFlag, targetNamespace)

	resultArgs := buildSchedulerArgs(testConfig)

	for i, arg := range resultArgs {
		if arg == testConfig.TargetPodName {
			assert.Equal(t, resultArgs[i-1], fmt.Sprintf("--%s", cmd.TargetPodName), "target pod arg should come after TargetPodName flag")
		}
		if arg == testConfig.TargetPodNamespace {
			assert.Contains(t, resultArgs[i-1], fmt.Sprintf("--%s", cmd.TargetPodNamespace), "target pod namespace arg should come after TargetPodNamespace flag")
		}
		if arg == testConfig.Namespace {
			assert.Contains(t, resultArgs[i-1], fmt.Sprintf("--%s", cmd.HomeNamespaceFlag), "home namespace arg should come after HomeNamespace flag")
		}
		if arg == username {
			assert.Contains(t, resultArgs[i-1], fmt.Sprintf("--%s", cmd.AuthUsername), "username arg should come after AuthUsername flag")
		}
		if arg == secret {
			assert.Contains(t, resultArgs[i-1], fmt.Sprintf("--%s", cmd.AuthSecretFlag), "secret arg should come after AuthSecret flag")
		}
	}

	assert.Contains(t, resultArgs, testConfig.TargetPodName)
	assert.Contains(t, resultArgs, testConfig.TargetPodNamespace)
	assert.Contains(t, resultArgs, "--debug")

	assert.Contains(t, resultArgs, username)
	assert.Contains(t, resultArgs, secret)
	assert.Contains(t, resultArgs, targetNamespace)
}
