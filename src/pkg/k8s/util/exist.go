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
	"context"
	"fmt"
	"github.com/suecodelabs/cnfuzz/src/pkg/logger"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NamespaceExists(l logger.Logger, clientset kubernetes.Interface, namespace string) (bool, error) {
	found, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if len(namespace) < 1 || err != nil || found == nil {
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return false, nil
			}
			l.V(logger.InfoLevel).Error(err, "error while checking if namespace exists")
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func PodExists(l logger.Logger, clientSet kubernetes.Interface, namespace string, podName string) (bool, error) {
	found, err := clientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	l.V(logger.DebugLevel).Info(fmt.Sprintf("found a pod while checking if pod %s exists", podName), "podName", podName)
	l.V(logger.PerformanceTestLevel).Info("found pod", "foundPod", found)
	if err != nil || found == nil {
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return false, nil
			}
			l.V(logger.InfoLevel).Error(err, "error while checking if pod exists")
			return false, err
		}
		return false, nil
	}
	l.V(logger.DebugLevel).Info("pod exists", "podName", podName)
	return true, nil
}
