package util

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"log"
	"testing"
)

func TestNamespaceExists(t *testing.T) {
	testNamespaces := []*v1.Namespace{
		createNamespace("my-ns1"),
		createNamespace("my-ns2"),
	}
	clientSet := fake.NewSimpleClientset()
	for _, namespace := range testNamespaces {
		_, err := clientSet.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
		if err != nil {
			log.Fatalln(fmt.Errorf("failed to create namespace in test kubernetes clientset: %w", err))
		}
	}
	result1 := NamespaceExist(clientSet, "my-ns1")
	assert.True(t, result1)
	result2 := NamespaceExist(clientSet, "ns-that-doesnt-exist")
	assert.False(t, result2)
}

func createNamespace(name string) *v1.Namespace {
	return &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}
