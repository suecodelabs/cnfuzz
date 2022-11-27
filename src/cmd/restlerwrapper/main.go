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

package main

import (
	"github.com/spf13/cobra"
	"github.com/suecodelabs/cnfuzz/src/internal/api_info"
	"github.com/suecodelabs/cnfuzz/src/internal/restler"
	"github.com/suecodelabs/cnfuzz/src/pkg/discovery/openapi"
	"github.com/suecodelabs/cnfuzz/src/pkg/logger"
	"log"
	"os"
)

type Command struct {
	command *cobra.Command
	*Args
}

type Args struct {
	isDebug         bool
	localConfig     bool
	dDocLoc         string
	targetPod       string
	targetNamespace string
	targetPort      int32
	dDocIp          string
	dryRun          bool
}

func main() {
	cmd := Command{
		command: &cobra.Command{
			Use: "rw <flags>",
		},
		Args: &Args{
			isDebug:         false,
			localConfig:     false,
			dDocLoc:         "/swagger/doc.json",
			targetPod:       "",
			targetNamespace: "default",
			targetPort:      0,
			dDocIp:          "",
			dryRun:          false,
		},
	}

	cmd.command.PersistentFlags().BoolVar(&cmd.Args.localConfig, "local-config", cmd.Args.localConfig, "Use the local kubeconfig instead of getting it from the cluster")
	cmd.command.PersistentFlags().StringVar(&cmd.Args.dDocLoc, "d-doc", cmd.Args.dDocLoc, "Uri of the discovery document (open API document)")
	cmd.command.PersistentFlags().StringVar(&cmd.Args.targetPod, "pod", cmd.Args.targetPod, "Set the pod name of the target")
	cmd.command.PersistentFlags().StringVar(&cmd.Args.targetNamespace, "ns", cmd.Args.targetPod, "Set the namespace of the target")
	cmd.command.PersistentFlags().Int32Var(&cmd.Args.targetPort, "port", cmd.Args.targetPort, "Set the port of the target service")
	cmd.command.PersistentFlags().BoolVarP(&cmd.Args.isDebug, "debug", "d", cmd.isDebug, "Enable debug mode")
	cmd.command.PersistentFlags().StringVar(&cmd.Args.dDocIp, "ddoc-ip", cmd.Args.dDocIp, "Dev flag: Overwrite the IP that is used to get the OpenApi doc")
	cmd.command.PersistentFlags().BoolVar(&cmd.dryRun, "dry-run", cmd.Args.dryRun, "Dev flag: Do a dry run, run without executing the Restler commands")

	cmd.command.Run = func(_ *cobra.Command, _ []string) {
		l := logger.CreateLogger(cmd.Args.isDebug)
		if cmd.Args.isDebug {
			l.V(logger.InfoLevel).Info("running in debug mode")
		}
		run(l, *cmd.Args)
	}

	if err := cmd.command.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func run(l logger.Logger, args Args) {
	var ports []int32
	if args.targetPort != 0 { // if ports is empty TryGetOpenApiDoc will guess the port
		ports = append(ports, args.targetPort)
	}
	l.V(logger.DebugLevel).Info("fetching info from target ...")
	info := api_info.CollectInfo(l, args.targetPod, args.targetNamespace, args.dDocIp, args.dDocLoc, ports, args.localConfig)
	if !args.dryRun {
		l.V(logger.DebugLevel).Info("writing OpenApi document to a file so Restler can pick it up later")
		writeDocToFile(l, info.UnparsedApiDoc)
	}

	l.V(logger.DebugLevel).Info("executing Restler commands")
	restler.ExecuteRestlerCmds(l, args.dryRun, info)
	l.V(logger.InfoLevel).Info("job finished, exiting now ...")
}

func writeDocToFile(l logger.Logger, apiDoc openapi.UnParsedOpenApiDoc) {
	b, err := apiDoc.DocFile.MarshalJSON()
	if err != nil {
		l.FatalError(err, "failed to marshal OpenApi doc to bytes")
	} else {
		err := os.Mkdir("/openapi", os.FileMode(0755))
		if err != nil {
			l.FatalError(err, "failed to create 'openapi' dir to write OpenApi doc into")
		}
		err = os.WriteFile("/openapi/doc.json", b, os.FileMode(0644))
		if err != nil {
			l.FatalError(err, "failed to write OpenApi doc to fs")
		}
	}
}
