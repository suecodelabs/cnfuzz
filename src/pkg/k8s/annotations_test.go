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
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"testing"
)

func TestGetAnnotations(t *testing.T) {
	ignoreMeVal := true
	fuzzMeVal := false
	oaDocVal := "/swagger/swagger.json"
	secretVal := "verysecret"
	unameVal := "me"
	testMeta := &metav1.ObjectMeta{
		Annotations: map[string]string{
			fmt.Sprintf("%s/%s", AnnotationPrefix, IgnoreMeAnno):   strconv.FormatBool(ignoreMeVal),
			fmt.Sprintf("%s/%s", AnnotationPrefix, FuzzMeAnno):     strconv.FormatBool(fuzzMeVal),
			fmt.Sprintf("%s/%s", AnnotationPrefix, OpenApiDocAnno): oaDocVal,
			fmt.Sprintf("%s/%s", AnnotationPrefix, SecretAnno):     secretVal,
			fmt.Sprintf("%s/%s", AnnotationPrefix, UsernameAnno):   unameVal,
		},
	}
	result := GetAnnotations(testMeta)
	assert.Equal(t, ignoreMeVal, result.IgnoreMe)
	assert.Equal(t, fuzzMeVal, result.FuzzMe)
	assert.Equal(t, oaDocVal, result.OpenApiDocLocation)
	assert.Equal(t, secretVal, result.Secret)
	assert.Equal(t, unameVal, result.Username)
}

func TestGetAnnotationsFromMeta(t *testing.T) {
	testAnno := "test"
	testValue := "myvalue"
	testMeta := &metav1.ObjectMeta{
		Annotations: map[string]string{
			fmt.Sprintf("%s/%s", AnnotationPrefix, testAnno): testValue,
		},
	}
	result := getAnnotationFromMeta(testMeta, testAnno)
	assert.Equal(t, testValue, result)
}
