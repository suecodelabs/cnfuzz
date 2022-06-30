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

package command

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/suecodelabs/cnfuzz/src/health"
	"github.com/suecodelabs/cnfuzz/src/kubernetes"
	"github.com/suecodelabs/cnfuzz/src/log"
	"github.com/suecodelabs/cnfuzz/src/persistence/repository"
)

type schedulerCmd struct {
	cmd *cobra.Command

	schedulerArgs
	BaseArgs
}

type schedulerArgs struct {
	cacheSolution string
	redisHostName string
	redisPort     string
}

func createScheduler() *schedulerCmd {
	s := &schedulerCmd{
		cmd: &cobra.Command{
			Use:   "scheduler",
			Short: "Watch cluster and schedule new fuzzing jobs.",
		},
	}

	s.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return s.Run()
	}

	s.cmd.PreRun = func(cmd *cobra.Command, args []string) {
		BasePreRun(s.BaseArgs)
	}

	s.cmd.PersistentFlags().StringVarP(&s.cacheSolution, "cache", "", "redis", "Select which caching solution to use (options: 'redis', 'in_memory'")
	s.cmd.PersistentFlags().StringVarP(&s.redisHostName, "redis-hostname", "", "redis-master", "The Redis hostname that the scheduler will use for caching purposes.")
	s.cmd.PersistentFlags().StringVarP(&s.redisPort, "redis-port", "", "6379", "The Redis port that the scheduler will use for caching purposes.")

	AddBaseFlags(s.cmd, &s.BaseArgs)

	return s
}

func (schedCmd schedulerCmd) Run() error {
	log.L().Info("running scheduler schedCmd")
	// Running as "scheduler" starting new jobs when new API's start

	healthChecker := health.NewChecker()

	// Init repositories for persistence
	// Storage is only necessary in this "scheduler" mode
	strCacheSolution := viper.GetString(schedCmd.cacheSolution)
	repoType, repoErr := repository.RepoTypeFromString(strCacheSolution)
	if repoErr != nil {
		log.L().Fatalf("%s is not a valid repo type: %+v", strCacheSolution, repoErr)
	}

	repos := repository.InitRepositories(repoType, &healthChecker)
	go health.Serv(healthChecker)

	err := kubernetes.StartInformers(repos)
	if err != nil {
		log.L().Fatal(err)
	}
	return nil
}

func init() {
	schedCmd := createScheduler()
	RootCmd.AddCommand(schedCmd.cmd)
}
