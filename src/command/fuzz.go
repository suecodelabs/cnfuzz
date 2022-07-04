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
	"github.com/suecodelabs/cnfuzz/src/health"
	"github.com/suecodelabs/cnfuzz/src/kubernetes"
	"github.com/suecodelabs/cnfuzz/src/log"
)

type fuzzCmd struct {
	cmd *cobra.Command

	fuzzArgs
	BaseArgs
}

type fuzzArgs struct {
	pod          string
	podNamespace string
}

func createFuzz() *fuzzCmd {
	s := &fuzzCmd{
		cmd: &cobra.Command{
			Use:   "fuzz",
			Short: "Fuzz a pod.",
		},
	}

	s.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return s.Run()
	}

	s.cmd.PreRun = func(cmd *cobra.Command, args []string) {
		BasePreRun(s.BaseArgs)
	}

	s.cmd.PersistentFlags().StringVarP(&s.pod, "pod", "", "", "Kubernetes pod to target for fuzzing")
	s.cmd.PersistentFlags().StringVarP(&s.podNamespace, "namespace", "", "", "Namespace of the target pod")

	AddBaseFlags(s.cmd, &s.BaseArgs)

	return s
}

func (cmd fuzzCmd) Run() error {
	healthChecker := health.NewChecker()
	go health.Serv(healthChecker)
	err := kubernetes.FuzzPodWithName(cmd.podNamespace, cmd.pod)
	if err != nil {
		// TODO handle error
		log.L().Fatal(err)
	}
	return nil
}

func init() {
	fuzzCmd := createFuzz()
	RootCmd.AddCommand(fuzzCmd.cmd)
}
