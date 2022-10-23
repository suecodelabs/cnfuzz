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

package restlerwrapper

import (
	"fmt"
	"github.com/suecodelabs/cnfuzz/src/pkg/logger"
	"github.com/suecodelabs/cnfuzz/src/pkg/restlerwrapper/auth"
)

func CreateRestlerCompileCommand(l logger.Logger) (cmd string, args []string) {
	cmd = "dotnet"
	args = []string{"/RESTler/restler/Restler.dll", "compile", "--api_spec", "/openapi/doc.json"}
	return cmd, args
}

// CreateRestlerCommand creates command string that can be run inside the RESTler container
// the command string consists of a compile command that analyzes the OpenAPI spec and generates a fuzzing grammar
// and the fuzz command itself
func CreateRestlerCommand(l logger.Logger, tokenSource auth.ITokenSource, targetIp, targetPort, targetScheme, timeBudget string) (cmd string, args []string) {
	l.V(logger.DebugLevel).Info(fmt.Sprintf("using %s:%s for restler", targetIp, targetPort), "targetIp", targetIp, "targetPort", targetPort)

	// Please, UNIX philosophy people.
	cmd = "dotnet"
	args = []string{"/RESTler/restler/Restler.dll", "fuzz", "--grammar_file", "/Compile/grammar.py", "--dictionary_file", "/Compile/dict.json",
		"--target_ip", targetIp, "--target_port", targetPort, "--time_budget", timeBudget}

	if targetScheme == "https" {
		l.V(logger.DebugLevel).Info("using SSL in Restler")
	} else {
		l.V(logger.DebugLevel).Info("not using SSL in Restler")
		args = append(args, "--no_ssl")
	}

	if tokenSource != nil {
		// create a new auth token using the tokensource
		tok, tokErr := tokenSource.Token()
		if tokErr != nil {
			l.V(logger.ImportantLevel).Error(tokErr, "error while getting a new auth token")
		} else {
			token := fmt.Sprintf("%s: %s", "Authorization", tok.CreateAuthHeaderValue(l))
			if tokErr == nil && len(token) > 0 {
				authCmd := fmt.Sprintf("\"python3 /scripts/auth.py '%s' '%s'\"", "FIX_ME", token)
				// Use a high refresh interval because we have a static token (for now?)
				args = append(args, "--token_refresh_interval", "999999", "--token_refresh_command", authCmd)
			}
		}
	}
	return cmd, args
}
