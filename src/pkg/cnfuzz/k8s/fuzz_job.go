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
	"github.com/suecodelabs/cnfuzz/src/pkg/cnfuzz/k8s/job"
	config "github.com/suecodelabs/cnfuzz/src/pkg/config"
	"github.com/suecodelabs/cnfuzz/src/pkg/discovery/openapi"
	"github.com/suecodelabs/cnfuzz/src/pkg/logger"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func StartFuzzJob(l logger.Logger, client kubernetes.Interface, cnfConfig *config.CnFuzzConfig, overwrites config.DDocOverwrites, pod *v1.Pod) error {
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

	// Check if the open api doc exists
	apiDesc, err := openapi.TryGetOpenApiDoc(l, ip, ports, oaLocs)
	if err != nil {
		return fmt.Errorf("error while retrieving OpenAPI document from target %s: %w", pod.Name, err)
	}

	restlerJob := job.CreateRestlerWrapperJob(l, pod, cnfConfig, apiDesc)
	createdJob, err := client.BatchV1().Jobs(restlerJob.Namespace).Create(context.TODO(), restlerJob, metav1.CreateOptions{})
	if err != nil {
		l.V(logger.ImportantLevel).Error(err, "error while starting restler job", "restlerJobName", restlerJob.Name, "restlerJobNamespace", restlerJob.Namespace, "targetName", pod.Name)
	}

	// TODO wait until the job is finished
	var _ = createdJob

	l.V(logger.InfoLevel).Info("completed job", "jobName", restlerJob.Name)

	return nil
}
