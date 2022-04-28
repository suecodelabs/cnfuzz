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

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	IsDebug = "debug"

	// start job with Kubernetes
	TargetPodName      = "pod"
	TargetPodNamespace = "namespace" // Namespace that target lives in

	// misc Kubernetes flags
	InsideClusterFlag  = "inside-cluster"
	OnlyFuzzMarkedFlag = "only-marked"
	SchedulerImageFlag = "scheduler-img"
	HomeNamespaceFlag  = "home-ns" // Namespace to start containers in (jobs etc.)

	// fuzzing related flags
	RestlerInitImageFlag   = "restler-init-img"
	RestlerImageFlag       = "restler-img"
	RestlerTelemetryOptOut = "restler-telemetry-opt-out"
	RestlerTimeBudget      = "restler-time-budget"
	RestlerCpuLimit        = "restler-cpu-limit"
	RestlerMemoryLimit     = "restler-memory-limit"
	RestlerCpuRequest      = "restler-cpu-request"
	RestlerMemoryRequest   = "restler-memory-request"

	// caching related flags
	CacheSolution = "cache"
	RedisHostName = "redis-hostname"
	RedisPort     = "redis-port"

	// auth related flags
	AuthUsername   = "username"
	AuthSecretFlag = "secret"

	//s3 flags
	S3EndpointUrlFlag = "s3-endpoint"
	S3ReportBucket    = "s3-bucket"
	S3AccessKey       = "s3-access"
	S3SecretKey       = "s3-secret"
)

// SetupFlags registers all flags with viper
func SetupFlags(rootCmd *cobra.Command) {
	// Debug flag
	rootCmd.Flags().BoolP(IsDebug, "d", false, "Debug mode")
	_ = viper.BindPFlag(IsDebug, rootCmd.Flags().Lookup(IsDebug))

	registerDirectFuzzingFlags(rootCmd)

	registerKubernetesFlags(rootCmd)

	registerCacheFlags(rootCmd)

	rootCmd.Flags().StringP(RestlerInitImageFlag, "", "curlimages/curl:7.81.0", "Init Image for preparing RESTler runtime")
	_ = viper.BindPFlag(RestlerInitImageFlag, rootCmd.Flags().Lookup(RestlerInitImageFlag))

	rootCmd.Flags().StringP(RestlerTelemetryOptOut, "", "", "Opt out for RESTler telemetry collection.")
	_ = viper.BindPFlag(RestlerTelemetryOptOut, rootCmd.Flags().Lookup(RestlerTelemetryOptOut))

	rootCmd.Flags().StringP(RestlerImageFlag, "", "mcr.microsoft.com/restlerfuzzer/restler:v7.4.0", "RESTler image to use (https://hub.docker.com/_/microsoft-restlerfuzzer-restler)")
	_ = viper.BindPFlag(RestlerImageFlag, rootCmd.Flags().Lookup(RestlerImageFlag))

	registerAuthFlags(rootCmd)
	registerS3Flags(rootCmd)
}

// registerDirectFuzzingFlags registers flags used when directly fuzzing a target
func registerDirectFuzzingFlags(rootCmd *cobra.Command) {
	rootCmd.Flags().StringP(TargetPodName, "", "", "Kubernetes pod to target for fuzzing")
	_ = viper.BindPFlag(TargetPodName, rootCmd.Flags().Lookup(TargetPodName))

	rootCmd.Flags().StringP(TargetPodNamespace, "n", "default", "Namespace for the target pod")
	_ = viper.BindPFlag(TargetPodNamespace, rootCmd.Flags().Lookup(TargetPodNamespace))

	rootCmd.Flags().StringP(RestlerTimeBudget, "", "1", "Maximum hours a Fuzzing Job may take.")
	_ = viper.BindPFlag(RestlerTimeBudget, rootCmd.Flags().Lookup(RestlerTimeBudget))

}

// registerKubernetesFlags registers flags for Kubernetes configuration
func registerKubernetesFlags(rootCmd *cobra.Command) {
	rootCmd.Flags().BoolP(InsideClusterFlag, "k", true, "Tells the fuzzer that it is running inside Kubernetes")
	_ = viper.BindPFlag(InsideClusterFlag, rootCmd.Flags().Lookup(InsideClusterFlag))

	rootCmd.Flags().BoolP(OnlyFuzzMarkedFlag, "m", false, "Only fuzz pods that have a '\"cnfuzz/enable\": \"true\"' annotation")
	_ = viper.BindPFlag(OnlyFuzzMarkedFlag, rootCmd.Flags().Lookup(OnlyFuzzMarkedFlag))

	// TODO change current temp default image to actual image ones it exists
	defaultJImg := ""
	rootCmd.PersistentFlags().StringP(SchedulerImageFlag, "i", defaultJImg, "Image used for the Kubernetes job, you can use this to change image version or replace the entire image")
	_ = viper.BindPFlag(SchedulerImageFlag, rootCmd.PersistentFlags().Lookup(SchedulerImageFlag))

	rootCmd.Flags().StringP(HomeNamespaceFlag, "", "default", "Namespace to start fuzzing containers in")
	_ = viper.BindPFlag(HomeNamespaceFlag, rootCmd.Flags().Lookup(HomeNamespaceFlag))

	rootCmd.Flags().Int64P(RestlerCpuLimit, "", 500, "Maximum amount of (milli) CPU a Fuzzing Job may use.")
	_ = viper.BindPFlag(RestlerCpuLimit, rootCmd.Flags().Lookup(RestlerCpuLimit))

	rootCmd.Flags().Int64P(RestlerMemoryLimit, "", 500, "Maximum memory (Mi) a Fuzzing Job may use.")
	_ = viper.BindPFlag(RestlerMemoryLimit, rootCmd.Flags().Lookup(RestlerMemoryLimit))

	rootCmd.Flags().Int64P(RestlerCpuRequest, "", 500, "Maximum amount of (milli) CPU a Fuzzing Job may request.")
	_ = viper.BindPFlag(RestlerCpuRequest, rootCmd.Flags().Lookup(RestlerCpuRequest))

	rootCmd.Flags().Int64P(RestlerMemoryRequest, "", 500, "Maximum memory (Mi) a Fuzzing Job may request.")
	_ = viper.BindPFlag(RestlerMemoryRequest, rootCmd.Flags().Lookup(RestlerMemoryRequest))
}

// registerAuthFlags registers flags for auth
func registerAuthFlags(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().StringP(AuthUsername, "", "fuzz-client", "Username to be used in auth (only works for basic auth)")
	rootCmd.PersistentFlags().StringP(AuthSecretFlag, "", "", "Secret to be used for authentication")
	_ = viper.BindPFlag(AuthUsername, rootCmd.PersistentFlags().Lookup(AuthUsername))
	_ = viper.BindPFlag(AuthSecretFlag, rootCmd.PersistentFlags().Lookup(AuthSecretFlag))
	// TODO add a scopes flag
}

func registerCacheFlags(rootCmd *cobra.Command) {
	rootCmd.Flags().StringP(CacheSolution, "c", "redis", "Select which caching solution to use (options: 'redis', 'in_memory'")
	_ = viper.BindPFlag(CacheSolution, rootCmd.Flags().Lookup(CacheSolution))

	rootCmd.Flags().StringP(RedisHostName, "", "redis-master", "The Redis hostname that the scheduler will use for caching purposes.")
	_ = viper.BindPFlag(RedisHostName, rootCmd.Flags().Lookup(RedisHostName))

	rootCmd.Flags().StringP(RedisPort, "", "6379", "The Redis port that the scheduler will use for caching purposes.")
	_ = viper.BindPFlag(RedisPort, rootCmd.Flags().Lookup(RedisPort))
}

func registerS3Flags(rootCmd *cobra.Command) {
	rootCmd.Flags().StringP(S3EndpointUrlFlag, "", "", "Optional endpoint url of your S3 bucket, example: 'http://my-minio-fs:9000'")
	_ = viper.BindPFlag(S3EndpointUrlFlag, rootCmd.Flags().Lookup(S3EndpointUrlFlag))

	rootCmd.Flags().StringP(S3ReportBucket, "s", "", "S3 bucket to copy fuzzing results to")
	_ = viper.BindPFlag(S3ReportBucket, rootCmd.Flags().Lookup(S3ReportBucket))

	rootCmd.Flags().StringP(S3AccessKey, "", "", "Access key of your S3 instance")
	_ = viper.BindPFlag(S3AccessKey, rootCmd.Flags().Lookup(S3AccessKey))

	rootCmd.Flags().StringP(S3SecretKey, "", "", "Secret key of your S3 instance")
	_ = viper.BindPFlag(S3SecretKey, rootCmd.Flags().Lookup(S3SecretKey))
}
