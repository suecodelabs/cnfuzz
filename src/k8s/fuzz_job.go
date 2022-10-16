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
	"context"
	"fmt"
	"github.com/suecodelabs/cnfuzz/src/auth"
	"github.com/suecodelabs/cnfuzz/src/config"
	"github.com/suecodelabs/cnfuzz/src/discovery/openapi"
	"github.com/suecodelabs/cnfuzz/src/k8s/job"
	"github.com/suecodelabs/cnfuzz/src/logger"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
)

func StartFuzzJobWithName(l logger.Logger, client kubernetes.Interface, cnfConfig *config.CnFuzzConfig, overwrites config.Overwrites, podName, podNamespace string) {
	pod, err := client.CoreV1().Pods(podNamespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		l.V(logger.ImportantLevel).Error(err, "error while getting pod info from cluster")
		os.Exit(1)
	}

	err = StartFuzzJob(l, client, cnfConfig, overwrites, pod)
	if err != nil {
		l.V(logger.ImportantLevel).Error(err, "error while fuzzing pod")
		os.Exit(1)
	}
}

func StartFuzzJob(l logger.Logger, client kubernetes.Interface, cnfConfig *config.CnFuzzConfig, overwrites config.Overwrites, pod *v1.Pod) error {
	annos := GetAnnotations(&pod.ObjectMeta)

	var ip string
	if len(overwrites.DiscoveryDocIP) > 0 {
		ip = overwrites.DiscoveryDocIP
	} else {
		ip = pod.Status.PodIP
	}

	var ports []int32
	if overwrites.DiscoveryDocPort != 0 {
		ports = append(ports, overwrites.DiscoveryDocPort)
	} else {
		for _, container := range pod.Spec.Containers {
			for _, port := range container.Ports {
				ports = append(ports, port.ContainerPort)
			}
		}
	}

	var oaLocs []string
	if len(annos.OpenApiDocLocation) > 0 {
		oaLocs = append(oaLocs, annos.OpenApiDocLocation)
	} else {
		oaLocs = openapi.GetCommonOpenApiLocations()
	}

	apiDesc, err := openapi.TryGetOpenApiDoc(l, ip, ports, oaLocs)
	if err != nil {
		return fmt.Errorf("error while retrieving OpenAPI document from target %s: %w", pod.Name, err)
	}

	// Tokensource can be nil !!! this means the API doesn't have any security (specified in the discovery doc ...)
	tokenSource, authErr := auth.CreateTokenSourceFromSchemas(l, apiDesc.SecuritySchemes, cnfConfig.AuthConfig.Username, cnfConfig.AuthConfig.Secret)
	if authErr != nil {
		l.V(logger.ImportantLevel).Error(authErr, "error while building auth token source")
		return authErr
	}

	restlerJob := job.CreateRestlerJob(pod, cnfConfig, apiDesc)
	j := restlerJob.CreateRestlerJob(l, tokenSource)
	createdJob, err := client.BatchV1().Jobs(j.Namespace).Create(context.TODO(), j, metav1.CreateOptions{})
	if err != nil {
		l.V(logger.ImportantLevel).Error(err, "error while starting restler job", "restlerJobName", j.Name, "restlerJobNamespace", j.Namespace, "targetName", pod.Name)
	}

	// TODO wait until the job is finished
	var _ = createdJob

	l.V(logger.InfoLevel).Info("completed job", "jobName", restlerJob.JobName)

	return nil
}
