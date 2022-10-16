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

package k8s

import (
	"github.com/suecodelabs/cnfuzz/src/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
)

// CreateClientset create a client to interact with the Kubernetes API
func CreateClientset(l logger.Logger, insideCluster bool) (clientset kubernetes.Interface) {
	var config *rest.Config
	var err error
	if insideCluster {
		// Inside cluster:
		l.V(logger.InfoLevel).Info("using config from cluster to create API client set")
		config, err = rest.InClusterConfig()
	} else {
		// Outside cluster:
		l.V(logger.InfoLevel).Info("using local kube config to create API client set")
		config, err = ctrl.GetConfig()
	}
	if err != nil {
		l.FatalError(err, "failed to get config for creating K8S API client set")
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		l.FatalError(err, "failed to create clientset from K8S config")
	}

	return clientset
}
