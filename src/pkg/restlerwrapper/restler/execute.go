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

package restler

import (
	"fmt"
	"github.com/suecodelabs/cnfuzz/src/pkg/logger"
	"github.com/suecodelabs/cnfuzz/src/pkg/restlerwrapper"
	"os/exec"
	"strings"
)

func ExecuteRestlerCmds(l logger.Logger, dryRun bool, info restlerwrapper.TargetInfo) {
	compileCmd, compileArgs := CreateRestlerCompileCommand(l)
	if !dryRun {
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

	restlerCmd, restlerArgs := CreateRestlerCommand(l, info.TokenSource, info.TargetAddr, info.ApiDesc.DiscoveryDoc.Port(), info.ApiDesc.Title, info.ApiDesc.DiscoveryDoc.Scheme, "1")
	if !dryRun {
		out, err := exec.Command(restlerCmd, restlerArgs...).Output()
		if err != nil {
			l.V(logger.InfoLevel).Info(fmt.Sprintf("restler output:\n%s", string(out[:])))
			l.FatalError(err, "error while executing restler fuzzing", "cmd_output", string(out[:]))
		}
		l.V(logger.DebugLevel).Info(string(out[:]))
	} else {
		fullCmd := restlerCmd + " " + strings.Join(restlerArgs, " ")
		l.V(logger.DebugLevel).Info("(running as dry run) generated restler cmd:")
		l.V(logger.DebugLevel).Info(fullCmd)
	}
}
