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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "cnfuzz targetUrl",
		Short: "Native Cloud Fuzzer is a fuzzer build for native cloud environments",
		Long: `Native Cloud Fuzzer is a fuzzer build for native cloud environments.
More info here:
https://github.com/suecodelabs/cnfuzz`,
		Args: cobra.NoArgs, // cobra.ExactArgs(1),
	}
)

// initializes the base command
func init() {
	cobra.OnInitialize(initConfig)

	SetupFlags(rootCmd)
}

// initConfig initializes viper configuration
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		cfgDir, err := os.UserConfigDir()
		cobra.CheckErr(err)

		cfgDir = cfgDir + "/ncfuzzer"
		viper.AddConfigPath(cfgDir)
	}
}

// Execute starts the base command
func Execute(runFunc func(cmd *cobra.Command, args []string)) error {
	rootCmd.Run = runFunc
	return rootCmd.Execute()
}
