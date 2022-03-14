package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/suecodelabs/cnfuzz/src/cmd"
	"github.com/suecodelabs/cnfuzz/src/kubernetes"
	"github.com/suecodelabs/cnfuzz/src/log"
	"github.com/suecodelabs/cnfuzz/src/persistence/repository"
	"github.com/suecodelabs/cnfuzz/src/serv"
)

func main() {
	if err := cmd.Execute(Run); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Run(command *cobra.Command, args []string) {
	isDebug := viper.GetBool(cmd.IsDebug)

	log.InitLogger(isDebug)
	logger := log.L()

	if isDebug {
		logger.Info("running in debug mode")
	}

	// Start the internal webserver
	go serv.Serv()

	//var targetUrl = args[0]
	//viper.Set(cmd.UrlArg, targetUrl)
	// url := viper.Get(cli.UrlArg)

	podName := viper.GetString(cmd.TargetPodName)
	namespace := viper.GetString(cmd.HomeNamespaceFlag)

	if podName != "" {
		// Fuzz this Kubernetes pod
		// All other information will be gathered through Kubernetes
		// Also IP information, so this probably won't work from outside the cluster
		logger.Infof("start fuzzing pod %s inside namespace %s\n", podName, namespace)

		err := kubernetes.FuzzPodWithName(namespace, podName)
		if err != nil {
			// TODO handle error
			logger.Fatal(err)
		}
	} else {
		// Running as "controller" starting new jobs when new API's start

		// Init repositories for persistence
		// Storage is only necessary in this "controller" mode
		repos := repository.InitRepositories()

		err := kubernetes.StartInformers(repos)
		if err != nil {
			logger.Fatal(err)
		}
	}
	logger.Debugf("exiting ...")
}
