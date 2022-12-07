package util

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"log"
	"testing"
	"time"
)

func TestWaitForPodReady(t *testing.T) {
	clientSet := fake.NewSimpleClientset()
	testNamespace := "test"
	testPod := createTestPod(clientSet, "test-pod", testNamespace)

	updateAfterDuration := time.Millisecond * 50 // millies
	timeoutAfter := time.Second * 5
	// Update the pod to ready after a duration
	go updatePodPhase(clientSet, testPod, updateAfterDuration)
	start := time.Now()
	err := WaitForPodReady(clientSet, context.TODO(), testPod.DeepCopy(), timeoutAfter)
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
	err := WaitForPodReady(clientSet, context.TODO(), testPod.DeepCopy(), time.Millisecond)
	assert.Error(t, err, "timed out waiting for the condition")
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
