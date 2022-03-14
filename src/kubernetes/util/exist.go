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
