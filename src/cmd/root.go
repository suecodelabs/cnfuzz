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

func init() {
	cobra.OnInitialize(initConfig)

	SetupFlags(rootCmd)
}

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

func Execute(runFunc func(cmd *cobra.Command, args []string)) error {
	rootCmd.Run = runFunc
	return rootCmd.Execute()
}
