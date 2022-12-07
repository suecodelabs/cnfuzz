package util

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"time"
)

func IsPodRunning(clientSet kubernetes.Interface, context context.Context, pod *corev1.Pod) wait.ConditionFunc {
	return func() (bool, error) {
		pod, err := clientSet.CoreV1().Pods(pod.Namespace).Get(context, pod.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if pod.Status.Phase == corev1.PodRunning {
			return true, nil
		}
		return false, nil
	}
}

func WaitForPodReady(clientSet kubernetes.Interface, context context.Context, pod *corev1.Pod, timeout time.Duration) error {
	return wait.PollImmediate(time.Second, timeout, IsPodRunning(clientSet, context, pod))
}
