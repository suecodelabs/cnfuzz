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

package command

import (
	"github.com/suecodelabs/cnfuzz/src/log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	RootCmd = &cobra.Command{
		Use:   "cnfuzz targetUrl",
		Short: "Native Cloud Fuzzer is a fuzzer build for native cloud environments",
		Long: `Native Cloud Fuzzer is a fuzzer build for native cloud environments.
More info here:
https://github.com/suecodelabs/cnfuzz`,
		// Args: cobra.NoArgs, // cobra.ExactArgs(1),
	}
)

type BaseArgs struct {
	isDebug        bool
	insideCluster  bool
	onlyMarked     bool
	schedulerImage string
	homeNamespace  string
}

func AddBaseFlags(cmd *cobra.Command, args *BaseArgs) {
	cmd.PersistentFlags().BoolVarP(&args.isDebug, "debug", "d", false, "Enable debug mode")
	cmd.PersistentFlags().Bool("inside-cluster", args.insideCluster, "Tells the fuzzer that it is running inside Kubernetes")
	cmd.PersistentFlags().Bool("only-marked", args.onlyMarked, "Only fuzz pods that have a '\"cnfuzz/enable\": \"true\"' annotation")
	cmd.PersistentFlags().String("scheduler-img", args.schedulerImage, "Image used for the Kubernetes job, you can use this to change image version or replace the entire image")
	cmd.PersistentFlags().String("home-ns", args.homeNamespace, "Namespace to start fuzzing containers in")
}

// initializes the base command
func init() {
	cobra.OnInitialize(initConfig)

	// SetupFlags(RootCmd)
}

// initConfig initializes viper configuration
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		cfgDir, err := os.UserConfigDir()
		cobra.CheckErr(err)

		cfgDir = cfgDir + "/cnfuzz"
		viper.AddConfigPath(cfgDir)
	}
}
func BasePreRun(args BaseArgs) {
	log.InitLogger(args.isDebug)
	logger := log.L()

	if args.isDebug {
		logger.Info("running in debug mode")
	}
}

// Execute starts the base command
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.L().Fatal(err)
	}
}
