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
		// Running as "scheduler" starting new jobs when new API's start

		// Init repositories for persistence
		// Storage is only necessary in this "scheduler" mode
		strCacheSolution := viper.GetString(cmd.CacheSolution)
		repoType, repoErr := repository.RepoTypeFromString(strCacheSolution)
		if repoErr != nil {
			log.L().Fatalf("%s is not a valid repo type: %+v", strCacheSolution, repoErr)
		}

		repos := repository.InitRepositories(repoType)

		err := kubernetes.StartInformers(repos)
		if err != nil {
			logger.Fatal(err)
		}
	}
	logger.Debugf("exiting ...")
}
