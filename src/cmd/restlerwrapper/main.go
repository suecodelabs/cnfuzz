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
	"github.com/suecodelabs/cnfuzz/src/pkg/discovery/openapi"
	"github.com/suecodelabs/cnfuzz/src/pkg/logger"
	"github.com/suecodelabs/cnfuzz/src/pkg/restlerwrapper"
	"log"
	"os/exec"
	"strings"
)

type Command struct {
	command *cobra.Command
	*Args
}

type Args struct {
	isDebug    bool
	dDocLoc    string
	targetIp   string
	targetPort int32
	dryRun     bool
}

func main() {
	cmd := Command{
		command: &cobra.Command{
			Use: "rw <flags>",
		},
		Args: &Args{
			isDebug:    false,
			dDocLoc:    "/swagger/doc.json",
			targetIp:   "",
			targetPort: 0,
			dryRun:     false,
		},
	}

	cmd.command.PersistentFlags().BoolVarP(&cmd.Args.isDebug, "debug", "d", cmd.isDebug, "Enable debug mode")
	cmd.command.PersistentFlags().StringVar(&cmd.Args.dDocLoc, "d-doc", cmd.Args.dDocLoc, "Uri of the discovery document (open API document)")
	cmd.command.PersistentFlags().StringVar(&cmd.Args.targetIp, "ip", cmd.Args.targetIp, "Set the IP of the target service")
	cmd.command.PersistentFlags().Int32Var(&cmd.Args.targetPort, "port", cmd.Args.targetPort, "Set the port of the target service")
	cmd.command.PersistentFlags().BoolVar(&cmd.dryRun, "dry-run", cmd.Args.dryRun, "Do a dry run, run without executing the Restler commands")

	cmd.command.Run = func(_ *cobra.Command, _ []string) {
		l := logger.CreateLogger(cmd.Args.isDebug)
		run(l, *cmd.Args)
	}

	if err := cmd.command.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func run(l logger.Logger, args Args) {
	// parse the passed arguments
	var ip string
	if len(args.targetIp) > 0 {
		ip = args.targetIp
	} else {
		l.Fatal("no target IP given")
	}

	var ports []int32
	if args.targetPort != 0 { // if ports is empty TryGetOpenApiDoc will guess the port
		ports = append(ports, args.targetPort)
	}

	openApiDocLocation := args.dDocLoc
	var oaLocs []string
	if len(openApiDocLocation) > 0 {
		oaLocs = append(oaLocs, openApiDocLocation)
	} else {
		oaLocs = openapi.GetCommonOpenApiLocations()
	}

	info := restlerwrapper.CollectInfoFromAddr(l, ip, ports, oaLocs, args.dryRun)

	compileCmd, compileArgs := restlerwrapper.CreateRestlerCompileCommand(l)
	if !args.dryRun {
		out, err := exec.Command(compileCmd, compileArgs...).Output()
		if err != nil {
			l.FatalError(err, "error while compiling restler resources")
		}
		l.V(logger.DebugLevel).Info(string(out[:]))
	} else {
		fullCmd := compileCmd + " " + strings.Join(compileArgs, " ")
		l.V(logger.DebugLevel).Info("(running as dry run) generated compile cmd:")
		l.V(logger.DebugLevel).Info(fullCmd)
	}

	restlerCmd, restlerArgs := restlerwrapper.CreateRestlerCommand(l, info.TokenSource, ip, info.ApiDesc.DiscoveryDoc.Port(), info.ApiDesc.DiscoveryDoc.Scheme, "1")
	if !args.dryRun {
		out, err := exec.Command(restlerCmd, restlerArgs...).Output()
		if err != nil {
			l.FatalError(err, "error while executing restler fuzzing")
		}
		l.V(logger.DebugLevel).Info(string(out[:]))
	} else {
		fullCmd := restlerCmd + " " + strings.Join(restlerArgs, " ")
		l.V(logger.DebugLevel).Info("(running as dry run) generated restler cmd:")
		l.V(logger.DebugLevel).Info(fullCmd)
	}
}
