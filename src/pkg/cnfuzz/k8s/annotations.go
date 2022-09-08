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

package k8s

import (
	"fmt"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	AnnotationPrefix = "cnfuzz"
	IgnoreMeAnno     = "ignore"
	FuzzMeAnno       = "enable"
	OpenApiDocAnno   = "open-api-doc"
	SecretAnno       = "secret"
	UsernameAnno     = "username"
)

// Annotations annotation values for annotations to be used inside Kubernetes configurations
type Annotations struct {
	IgnoreMe           bool
	FuzzMe             bool
	OpenApiDocLocation string
	Secret             string
	Username           string
}

// GetAnnotations gather annotations inside the metadata of a Kubernetes object
func GetAnnotations(objectMeta *metav1.ObjectMeta) Annotations {
	strIgnoreMe := getAnnotationFromMeta(objectMeta, IgnoreMeAnno)
	strFuzzMe := getAnnotationFromMeta(objectMeta, FuzzMeAnno)
	oaDocLoc := getAnnotationFromMeta(objectMeta, OpenApiDocAnno)
	secret := getAnnotationFromMeta(objectMeta, SecretAnno)
	username := getAnnotationFromMeta(objectMeta, UsernameAnno)

	ignoreMe, err := strconv.ParseBool(strIgnoreMe)
	if err != nil {
		// The value is invalid
		// just ignore the annotation
		ignoreMe = false
	}

	fuzzMe, err := strconv.ParseBool(strFuzzMe)

	return Annotations{
		IgnoreMe:           ignoreMe,
		FuzzMe:             fuzzMe,
		OpenApiDocLocation: oaDocLoc,
		Secret:             secret,
		Username:           username,
	}

}

// getAnnotationFromMeta get a single annotation value from Kubernetes object meta
func getAnnotationFromMeta(objectMeta *metav1.ObjectMeta, annotationName string) string {
	return objectMeta.Annotations[fmt.Sprintf("%s/%s", AnnotationPrefix, annotationName)]
}
