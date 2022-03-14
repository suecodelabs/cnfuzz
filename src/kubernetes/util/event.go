package util

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func IsKubernetesEvent(object *metav1.ObjectMeta) bool {
	return object.Name == "kubernetes" || object.Namespace == "kube-system" || object.Namespace == "kube-node-lease" || object.Namespace == "kube-pulbic"
}

func IsFuzzerEvent(object *metav1.ObjectMeta) bool {
	return strings.HasPrefix(object.Name, "cnfuzz")
}
