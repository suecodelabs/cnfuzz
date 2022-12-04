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

package persistence

import "fmt"

type StorageType int

const (
	Redis StorageType = iota
	InMemory
)

// StorageTypes that cnfuzz supports.
// Currently contains Redis and InMemory support.
// These strings are used to map the value from the config to a StorageType type.
var StorageTypes = [2]string{"redis", "in_memory"}

// String() returns the string equivalent of the enumeration
func (s StorageType) String() string {
	return StorageTypes[s]
}

// StorageTypeFromString creates a StorageType from a string.
// For example, calling this method with 'redis' returns Redis.
func StorageTypeFromString(value string) (StorageType, error) {
	for i := 0; i < len(StorageTypes); i++ {
		if StorageTypes[i] == value {
			return StorageType(i), nil
		}
	}
	return StorageType(0), fmt.Errorf("unkown storage type")
}
