// Copyright 2022 Sue B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubernetes

import (
	"fmt"

	"github.com/suecodelabs/cnfuzz/src/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
)

// CreateClientSet create a client to interact with the Kubernetes API
func CreateClientSet(insideCluster bool) (clientset kubernetes.Interface, err error) {
	logger := log.L()
	if insideCluster {
		// Inside cluster:
		logger.Debug("using in cluster config to build Kubernetes client set")
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("error while getting Kubernetes config from inside the cluster: %w", err)
		}
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create a clientset out of the received Kubernetes config: %w", err)
		}
		return clientset, nil
	} else {
		// Outside cluster:
		logger.Debugf("using local config to build Kubernetes client set")
		var err error
		clientset, err = kubernetes.NewForConfig(ctrl.GetConfigOrDie())
		if err != nil {
			return nil, fmt.Errorf("failed to get Kubernetes config from local machine: %w", err)
		}
		return clientset, nil
	}
}
