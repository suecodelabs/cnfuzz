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
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"github.com/suecodelabs/cnfuzz/src/pkg/cnfuzz/config"
	"github.com/suecodelabs/cnfuzz/src/pkg/discovery/openapi"
	"github.com/suecodelabs/cnfuzz/src/pkg/logger"
	"github.com/suecodelabs/cnfuzz/src/pkg/restlerwrapper"
	"github.com/suecodelabs/cnfuzz/src/pkg/restlerwrapper/auth"
	"log"
	"os"
	"os/exec"
)

type Command struct {
	command *cobra.Command
	*Args
}

type Args struct {
	isDebug     bool
	localConfig bool
	configFile  string
	printConfig bool
}

func main() {
	cmd := Command{
		command: &cobra.Command{
			Use: "rw <flags>",
		},
		Args: &Args{
			isDebug:     false,
			localConfig: false,
			configFile:  "/config/config.yaml",
			printConfig: false,
		},
	}

	cmd.command.PersistentFlags().BoolVarP(&cmd.Args.isDebug, "debug", "d", cmd.isDebug, "Enable debug mode")
	cmd.command.PersistentFlags().BoolVar(&cmd.Args.localConfig, "local-config", cmd.Args.localConfig, "Use the local kubeconfig instead of getting it from the cluster")
	cmd.command.PersistentFlags().BoolVar(&cmd.Args.printConfig, "print-config", cmd.Args.printConfig, "Print the config file")
	cmd.command.PersistentFlags().StringVar(&cmd.Args.configFile, "config", cmd.Args.configFile, "Location of the config file to use")

	cmd.command.Run = func(_ *cobra.Command, _ []string) {
		l := logger.CreateLogger(cmd.Args.isDebug)
		run(l, *cmd.Args)
	}

	if err := cmd.command.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func run(l logr.Logger, args Args) {
	cnf := config.LoadConfigFile(l, args.configFile, args.printConfig)

	var _ = cnf

	podName := "somepodname"
	ip := "0.0.0.0"
	ports := []int32{8080}
	openApiDocLocation := ""

	var oaLocs []string
	if len(openApiDocLocation) > 0 {
		oaLocs = append(oaLocs, openApiDocLocation)
	} else {
		oaLocs = openapi.GetCommonOpenApiLocations()
	}

	apiDesc, err := openapi.TryGetOpenApiDoc(l, ip, ports, oaLocs)
	if err != nil {
		l.V(logger.ImportantLevel).Error(err, "error while retrieving OpenAPI document")
		os.Exit(1)
	}

	// Tokensource can be nil !!! this means the API doesn't have any security (specified in the discovery doc ...)
	tokenSource, authErr := auth.CreateTokenSourceFromSchemas(l, apiDesc.SecuritySchemes, cnf.AuthConfig.Username, cnf.AuthConfig.Secret)
	if authErr != nil {
		l.V(logger.ImportantLevel).Error(authErr, "error while building auth token source")
		os.Exit(1)
	}
	compileCmd, compileArgs := restlerwrapper.CreateRestlerCompileCommand(l)
	out, err := exec.Command(compileCmd, compileArgs...).Output()
	if err != nil {
		l.V(logger.ImportantLevel).Error(err, "error while compiling restler resources")
		os.Exit(1)
	}
	l.V(logger.DebugLevel).Info(string(out[:]))

	restlerCmd, restlerArgs := restlerwrapper.CreateRestlerCommand(l, tokenSource, ip, apiDesc.DiscoveryDoc.Port(), apiDesc.DiscoveryDoc.Scheme, "1", podName)
	out, err = exec.Command(restlerCmd, restlerArgs...).Output()
	if err != nil {
		l.V(logger.ImportantLevel).Error(err, "error while executing restler fuzzing")
		os.Exit(1)
	}
	l.V(logger.DebugLevel).Info(string(out[:]))
}
