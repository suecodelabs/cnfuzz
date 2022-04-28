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

package job

import (
	"context"

	"github.com/suecodelabs/cnfuzz/src/auth"
	"github.com/suecodelabs/cnfuzz/src/config"
	"github.com/suecodelabs/cnfuzz/src/log"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// LaunchFuzzJob launches a job that retrieves the OpenAPI doc and kicks of a RESTler job
func LaunchFuzzJob(clientset kubernetes.Interface, kubeConfig *config.SchedulerConfig) (createdJob *batchv1.Job, err error) {
	jobs := clientset.BatchV1().Jobs(kubeConfig.Namespace)

	jobSpec := createSchedulerJob(kubeConfig)

	result, err := jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	return result, err
}

// LaunchRestlerJob starts a job with the RESTler fuzzer
func LaunchRestlerJob(clientset kubernetes.Interface, restlerConfig *config.FuzzerConfig, tokenSource auth.ITokenSource) (restlerJob *batchv1.Job, err error) {
	// Start the RESTler container as a job
	restlerJobSpec := createRestlerJob(restlerConfig, tokenSource)
	restlerJob, jErr := clientset.BatchV1().Jobs(restlerConfig.Namespace).Create(context.TODO(), restlerJobSpec, metav1.CreateOptions{})
	if jErr != nil {
		log.L().Errorf("failed to create restler job: %+v", jErr)
		return nil, jErr
	}
	return restlerJobSpec, nil
}
