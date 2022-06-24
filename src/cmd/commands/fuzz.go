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

type fuzzCmd struct {
	cmd *cobra.Command

	fuzzBuilderCommon
}

type fuzzBuilderCommon struct {
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

	s.cmd.PersistentFlags().StringVarP(&s.pod, "pod", "", "", "Kubernetes pod to target for fuzzing")
	s.cmd.PersistentFlags().StringVarP(&s.podNamespace, "namespace", "", "", "Namespace of the target pod")

	return s
}

func (cmd fuzzCmd) Run() error {
	return nil
}
