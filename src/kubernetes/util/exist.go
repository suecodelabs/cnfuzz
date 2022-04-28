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
	"context"
	"strings"

	"github.com/suecodelabs/cnfuzz/src/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NamespaceExist(clientset kubernetes.Interface, namespace string) bool {
	found, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if len(namespace) < 1 || err != nil || found == nil {
		if err != nil {
			log.L().Errorf("error while getting namespace to check if it exists: %+v", err)
		}
		return false
	}
	return true
}

func PodExist(clientSet kubernetes.Interface, namespace string, podName string) bool {
	found, err := clientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil || found == nil {
		log.L().Debugf("checking if pod exists:\nfound pod: %s\nerr: %+v", found, err)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return false
			}
			log.L().Errorf("error while getting pod to check if it exists: %+v", err)
		}
		return false
	} else {
		log.L().Debugf("found pod exists: %s", found)
	}
	return true
}
