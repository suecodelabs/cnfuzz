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
	"github.com/stretchr/testify/assert"
	"github.com/suecodelabs/cnfuzz/src/logger"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"log"
	"testing"
	"time"
)

func TestIsPodRunning(t *testing.T) {
	clientSet := fake.NewSimpleClientset()
	testNamespace := "test"
	testPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: testNamespace,
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodPending,
		},
	}
	podCopy := testPod.DeepCopy()
	isRunningFunc := IsPodRunning(clientSet, podCopy)

	_, err := clientSet.CoreV1().Pods(testNamespace).Create(context.TODO(), testPod, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to create test pod: %w", err))
	}
	result, err := isRunningFunc()
	if assert.NoError(t, err) {
		assert.False(t, result)
	}

	testPod.Status.Phase = corev1.PodRunning
	_, err = clientSet.CoreV1().Pods(testNamespace).Update(context.TODO(), testPod, metav1.UpdateOptions{})
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to update test pod to a running status: %w", err))
	}

	result, err = isRunningFunc()
	if assert.NoError(t, err) {
		assert.True(t, result)
	}
}

func TestWaitForPodReady(t *testing.T) {
	clientSet := fake.NewSimpleClientset()
	testNamespace := "test"
	testPod := createTestPod(clientSet, "test-pod", testNamespace)

	updateAfterDuration := time.Millisecond * 50 // millies
	timeoutAfter := time.Second * 5
	// Update the pod to ready after a duration
	go updatePodPhase(clientSet, testPod, updateAfterDuration)
	start := time.Now()
	err := WaitForPodReady(clientSet, testPod.DeepCopy(), timeoutAfter, time.Millisecond)
	end := time.Now()
	totalDuration := end.UnixMilli() - start.UnixMilli()
	if assert.NoError(t, err) {
		// Function should only return after the pod is set to ready
		assert.Greater(t, totalDuration, updateAfterDuration.Milliseconds())
		// Function should return before timeout
		assert.Less(t, totalDuration, timeoutAfter.Milliseconds())
	}
}

func TestWaitForPodReadyTimeout(t *testing.T) {
	clientSet := fake.NewSimpleClientset()
	testNamespace := "test"
	testPod := createTestPod(clientSet, "test-pod", testNamespace)
	err := WaitForPodReady(clientSet, testPod.DeepCopy(), time.Millisecond, time.Millisecond)
	assert.Error(t, err, "timed out waiting for the condition")
}

func TestIsJobReady(t *testing.T) {
	clientSet := fake.NewSimpleClientset()
	testJobName := "test-pod"
	testNamespace := "test"
	testJob := &v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testJobName,
			Namespace: testNamespace,
		},
	}
	_, _ = clientSet.BatchV1().Jobs(testNamespace).Create(context.TODO(), testJob, metav1.CreateOptions{})
	isReadyResult := IsJobReady(clientSet, testJobName, testNamespace)
	result, _ := isReadyResult()
	assert.False(t, result)

	testJob.Status.Conditions = append(testJob.Status.Conditions, v1.JobCondition{Type: v1.JobComplete})
	_, _ = clientSet.BatchV1().Jobs(testNamespace).Update(context.TODO(), testJob, metav1.UpdateOptions{})
	isReadyResult = IsJobReady(clientSet, testJobName, testNamespace)
	result, _ = isReadyResult()
	assert.True(t, result)
}

func TestWaitForJobReadyTimeout(t *testing.T) {
	clientSet := fake.NewSimpleClientset()
	testJobName := "test-job"
	testNamespace := "test-namespace"
	_ = createTestJob(clientSet, testJobName, testNamespace)
	err := WaitForJobReady(clientSet, testJobName, testNamespace, time.Millisecond, time.Millisecond)
	assert.Error(t, err, "timed out waiting for the condition")
}

func TestWaitForJobReady(t *testing.T) {
	l := logger.CreateDebugLogger()
	clientSet := fake.NewSimpleClientset()
	testJobName := "test-job"
	testNamespace := "test-namespace"
	testJob := createTestJob(clientSet, testJobName, testNamespace)

	updateAfterDuration := time.Millisecond * 50 // millies
	// Code polls api every 5 seconds, so timeout should be at least 6 seconds to do a proper test
	timeoutAfter := time.Second * 6

	// Update the job to complete after a duration
	testJob.Status.Conditions = append(testJob.Status.Conditions, v1.JobCondition{Type: v1.JobComplete})
	go updateJob(clientSet, testJob, updateAfterDuration)

	start := time.Now()
	err := WaitForJobReady(clientSet, testJobName, testNamespace, timeoutAfter, time.Millisecond)
	end := time.Now()
	totalDuration := end.UnixMilli() - start.UnixMilli()
	l.Info(fmt.Sprintf("test finished after %v", totalDuration))
	if assert.NoError(t, err) {
		// Function should only return after the job is set to complete
		assert.Greater(t, totalDuration, updateAfterDuration.Milliseconds())
		// Function should return before timeout
		assert.Less(t, totalDuration, timeoutAfter.Milliseconds())
	}
}

func createTestPod(clientSet kubernetes.Interface, podName string, namespace string) *corev1.Pod {
	testPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodPending,
		},
	}
	_, err := clientSet.CoreV1().Pods(namespace).Create(context.TODO(), testPod, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to create test pod: %w", err))
	}
	return testPod
}

func updatePodPhase(clientSet kubernetes.Interface, pod *corev1.Pod, waitDuration time.Duration) {
	time.Sleep(waitDuration)
	pod.Status.Phase = corev1.PodRunning
	_, err := clientSet.CoreV1().Pods(pod.Namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to update test pod: %w", err))
	}
}

func createTestJob(clientset kubernetes.Interface, jobName string, namespace string) *v1.Job {
	testJob := &v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
		},
		Status: v1.JobStatus{},
	}
	_, _ = clientset.BatchV1().Jobs(namespace).Create(context.TODO(), testJob, metav1.CreateOptions{})
	return testJob
}

func updateJob(clientset kubernetes.Interface, job *v1.Job, waitDuration time.Duration) {
	time.Sleep(waitDuration)
	_, err := clientset.BatchV1().Jobs(job.Namespace).Update(context.TODO(), job, metav1.UpdateOptions{})
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to update test pod: %w", err))
	}
}
