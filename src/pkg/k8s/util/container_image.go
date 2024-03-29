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

package util

import (
	"github.com/suecodelabs/cnfuzz/src/pkg/logger"
	"strings"
)

func GetImageName(fullName string) (name string, tags []string) {
	// format: registry/image:tag1:tag2

	if strings.ContainsAny(fullName, "/") {
		// If the registry contains a :, we need ignore it when getting the tags
		// This could happen when the registry contains a :port
		registrySepIndex := strings.LastIndex(fullName, "/")
		registryName := fullName[:registrySepIndex-1]
		if strings.ContainsAny(registryName, ":") {
			// Parsing the tags normally will cause problems because the registry name/ address contains a :
			tags = strings.Split(fullName, ":")
			nameEnd := len(tags[0]) + len(tags[1]) + 1
			imageName := fullName[:nameEnd]
			return imageName, tags[2:]
		}
		// Else the parsing will go as expected
	}
	tags = strings.Split(fullName, ":")
	name = tags[0]
	// Remove first item (image name)
	// tags = append(tags[:0], tags[1:]...)
	return name, tags[1:]
}

func SplitImageId(l logger.Logger, imageId string) (hash string, hashType string) {
	l.V(logger.DebugLevel).Info("splitting imageId into hash and a hash type", "imageId", imageId)
	// Format:
	// docker-pullable://localhost:5000/imagename@sha256:5add8f7cf10b367af0fd4d9779a48891d9083ab56a691065421571b4d4cf4789
	// sha256:5add8f7cf10b367af0fd4d9779a48891d9083ab56a691065421571b4d4cf4789
	if strings.Contains(imageId, ":") {
		hashSplit := strings.Split(imageId, ":")
		if len(hashSplit) > 2 {
			// Last piece is the hash
			hash = hashSplit[len(hashSplit)-1]

			// Split the rest and start at the second to last piece
			imageIdPart := hashSplit[len(hashSplit)-2]
			hashTypeSplit := strings.Split(imageIdPart, "@")
			// Last piece should be the hash type
			hashType = hashTypeSplit[len(hashTypeSplit)-1]
		} else if len(hashSplit) == 2 {
			if strings.Contains(hashSplit[0], "@") {
				imgSplit := strings.Split(hashSplit[0], "@")
				return hashSplit[1], imgSplit[len(hashSplit)-1]
			}
			return hashSplit[1], hashSplit[0]
		}
	}

	return hash, hashType
}
