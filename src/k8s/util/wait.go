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
	"time"

	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	typev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func IsJobReady(clientset kubernetes.Interface, jobName string, namespace string) wait.ConditionFunc {
	return func() (bool, error) {
		job, getErr := clientset.BatchV1().Jobs(namespace).Get(context.TODO(), jobName, metav1.GetOptions{})
		if getErr != nil {
			return false, getErr
		}
		if len(job.Status.Conditions) == 0 {
			return false, nil
		}
		for _, condition := range job.Status.Conditions {
			if condition.Type == v1.JobSuspended {
				return false, nil
			}
			// Maybe also need to check for failed conditions?
			// Some jobs are only completed when a jobs succeeds
		}
		return true, nil
	}
}

func IsPodRunning(clientSet kubernetes.Interface, pod *corev1.Pod) wait.ConditionFunc {
	return func() (bool, error) {
		pod, err := clientSet.CoreV1().Pods(pod.Namespace).Get(context.TODO(), pod.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if pod.Status.Phase == corev1.PodRunning {
			return true, nil
		}
		return false, nil
	}
}

func WaitForJobReady(clientset kubernetes.Interface, jobName string, namespace string, timeout time.Duration) error {
	// Poll every 5 seconds if job has completed
	return wait.PollImmediate(time.Second*5, timeout, IsJobReady(clientset, jobName, namespace))
}

func WaitForPodReady(clientSet kubernetes.Interface, pod *corev1.Pod, timeout time.Duration) error {
	return wait.PollImmediate(time.Second, timeout, IsPodRunning(clientSet, pod))
}

func GetPodsForSvc(svc *corev1.Service, namespace string, k8sClient typev1.CoreV1Interface) (*corev1.PodList, error) {
	set := labels.Set(svc.Spec.Selector)
	listOptions := metav1.ListOptions{LabelSelector: set.AsSelector().String()}
	pods, err := k8sClient.Pods(namespace).List(context.Background(), listOptions)
	return pods, err
}
