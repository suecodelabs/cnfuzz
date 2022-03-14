package config

import (
	"github.com/spf13/viper"
	"github.com/suecodelabs/cnfuzz/src/cmd"
	"github.com/suecodelabs/cnfuzz/src/discovery"
)

type FuzzConfig struct {
	ApiDescription   *discovery.WebApiDescription
	ClientId         string
	Secret           string
	KubernetesConfig *KubernetesFuzzConfig
	TimeBudget       string
}

func NewFuzzConfigWKubeConfig(apiDescription *discovery.WebApiDescription, kubeConfig *KubernetesFuzzConfig) *FuzzConfig {
	clientId := viper.GetString(cmd.AuthUsername)
	secret := viper.GetString(cmd.AuthSecretFlag)
	timeBudget := viper.GetString(cmd.RestlerTimeBudget)

	return &FuzzConfig{
		ApiDescription:   apiDescription,
		ClientId:         clientId,
		Secret:           secret,
		KubernetesConfig: kubeConfig,
		TimeBudget:       timeBudget,
	}
}

func NewFuzzConfig(apiDescription *discovery.WebApiDescription, podName string, namespace string) *FuzzConfig {
	knConfig := CreateKubernetesConfig(podName, namespace)
	return NewFuzzConfigWKubeConfig(apiDescription, knConfig)
}
