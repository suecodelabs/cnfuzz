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

package util

import (
	"bufio"
	"os"

	"github.com/suecodelabs/cnfuzz/src/log"
)

func WriteToFile(bytesToWrite *[]byte, filePath string) {
	// os.Mkdir("./reports", os.FileMode(0666))
	file, err := os.Create(filePath)
	if err != nil {
		log.L().Errorf("error while creating \"%s\" on filesystem: %+v", filePath, err)
		return
	}
	writer := bufio.NewWriter(file)
	_, err = writer.Write(*bytesToWrite)
	if err != nil {
		log.L().Errorf("error while writing a report to filesystem: %+v", err)
		return
	}
	writer.Flush()
}
