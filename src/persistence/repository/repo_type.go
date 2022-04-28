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

package repository

import "fmt"

type RepoType int

const (
	Redis RepoType = iota
	InMemory
)

var RepoTypes = [2]string{"redis", "in_memory"}

func (s RepoType) String() string {
	return RepoTypes[s]
}

func RepoTypeFromString(value string) (RepoType, error) {
	for i := 0; i < len(RepoTypes); i++ {
		if RepoTypes[i] == value {
			return RepoType(i), nil
		}
	}
	return RepoType(0), fmt.Errorf("unkown repo type")
}
