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
	"github.com/suecodelabs/cnfuzz/src/pkg/logger"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"log"
	"testing"
)

func TestNamespaceExists(t *testing.T) {
	l := logger.CreateDebugLogger()
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
	result1, err1 := NamespaceExists(l, clientSet, "my-ns1")
	assert.NoError(t, err1)
	assert.True(t, result1)
	result2, err2 := NamespaceExists(l, clientSet, "ns-that-doesnt-exist")
	assert.NoError(t, err2)
	assert.False(t, result2)
}

func createNamespace(name string) *v1.Namespace {
	return &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

// TODO test PodExists method
