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

package config

import (
	"github.com/suecodelabs/cnfuzz/src/pkg/logger"
	"os"
	"path/filepath"
)

func loadFile(l logger.Logger, configFile string, printFile bool) (*[]byte, error) {
	if configFile == "" {
		configFile = filepath.Join()
		if !pathExists(configFile) {
			// return nil
		}
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		// l.FatalError(err, "failed to read config file", "configFile", configFile)
		return nil, err
	}

	if printFile {
		l.V(logger.DebugLevel).Info("\n" + string(data[:]))
	}

	return &data, nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
