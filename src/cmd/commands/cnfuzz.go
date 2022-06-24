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

package commands

import (
	"github.com/spf13/cobra"
)

const (
	// misc Kubernetes flags
	InsideClusterFlag  = "inside-cluster"
	OnlyFuzzMarkedFlag = "only-marked"
	SchedulerImageFlag = "scheduler-img"
	HomeNamespaceFlag  = "home-ns" // Namespace to start containers in (jobs etc.)
)

type cnfuzzCmd struct {
	cmd *cobra.Command

	cnfuzzBuilderCommon
}

type cnfuzzBuilderCommon struct {
	debugMode      bool
	insideCluster  bool
	onlyMarked     bool
	schedulerImage string
	homeNamespace  string
}

func newCnFuzzCommand() {
	cc := &cnfuzzCmd{}

	rootCmd := &cobra.Command{
		Use:   "cnfuzz",
		Short: "Native Cloud Fuzzer is a fuzzer build for native cloud environments",
		Long: `Native Cloud Fuzzer is a fuzzer build for native cloud environments.
More info here:
https://github.com/suecodelabs/cnfuzz`,
		Args: cobra.NoArgs,
	}

	// Debug flag
	rootCmd.PersistentFlags().Bool("debug", cc.debugMode, "Debug mode")

	rootCmd.PersistentFlags().Bool("inside-cluster", cc.insideCluster, "Tells the fuzzer that it is running inside Kubernetes")
	rootCmd.PersistentFlags().Bool("only-marked", cc.onlyMarked, "Only fuzz pods that have a '\"cnfuzz/enable\": \"true\"' annotation")
	rootCmd.PersistentFlags().String("scheduler-img", cc.schedulerImage, "Image used for the Kubernetes job, you can use this to change image version or replace the entire image")
	rootCmd.PersistentFlags().String("home-ns", cc.homeNamespace, "Namespace to start fuzzing containers in")

	cc.cmd = rootCmd
}
