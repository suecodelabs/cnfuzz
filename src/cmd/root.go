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
		Short: "cnfuzz is a fuzzer build for Cloud Native environments",
		Long: `cnfuzz is a fuzzer build for Cloud Native environments.
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

		cfgDir = cfgDir + "/cnfuzz"
		viper.AddConfigPath(cfgDir)
	}
}

// Execute starts the base command
func Execute(runFunc func(cmd *cobra.Command, args []string)) error {
	rootCmd.Run = runFunc
	return rootCmd.Execute()
}
