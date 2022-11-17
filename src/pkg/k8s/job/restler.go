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

package job

import (
	"fmt"
	"github.com/suecodelabs/cnfuzz/src/pkg/config"
	"github.com/suecodelabs/cnfuzz/src/pkg/discovery/openapi"
	"github.com/suecodelabs/cnfuzz/src/pkg/logger"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO update this desc:

// CreateRestlerWrapperJob creates a Kubernetes Job for the cnfuzz wrapper around the RESTler fuzzer
// this includes an init container that gets the OpenAPI doc from the target API with curl and volumes for transferring the information
// it uses values from the FuzzConfig to configure the fuzz command that runs inside the RESTler container
// returned job hasn't started yet
func CreateRestlerWrapperJob(l logger.Logger, targetPod *v1.Pod, cnf *config.CnFuzzConfig, dDoc openapi.UnParsedOpenApiDoc) *batchv1.Job {
	restlerCnf := cnf.RestlerWrapperConfig.RestlerConfig
	imgCnf := cnf.RestlerWrapperConfig.ImageConfig

	jobName := "cnfuzz-job-" + targetPod.Name
	namespace := targetPod.Namespace
	containerName := "cnfuzz-job-" + targetPod.Name + "-restler"
	serviceAcc := cnf.RestlerWrapperConfig.ServiceAccount

	var containerImage string
	if len(imgCnf.Tag) > 0 {
		containerImage = fmt.Sprintf("%s:%s", imgCnf.Image, imgCnf.Tag)
	} else {
		containerImage = imgCnf.Image
	}
	pullPolicy := v1.PullIfNotPresent
	if len(imgCnf.PullPolicy) > 0 {
		pullPolicy = v1.PullPolicy(imgCnf.PullPolicy)
	}
	// targetIp := targetPod.Status.PodIP
	targetPort := dDoc.Uri.Port()
	targetDiscDocLoc := dDoc.Uri.Path // TODO Is this correct?
	telemetryOptOut := restlerCnf.TelemetryOptOut
	cpuRequest := resource.MustParse(restlerCnf.CpuRequest)
	memoryRequest := resource.MustParse(restlerCnf.MemoryRequest)
	cpuLimit := resource.MustParse(restlerCnf.CpuLimit)
	memoryLimit := resource.MustParse(restlerCnf.MemoryLimit)

	restlerWrapperArgs := []string{"--pod", targetPod.Name, "--ns", targetPod.Namespace, "--port", targetPort, "--d-doc", targetDiscDocLoc}
	debugMode := false // TODO take value from arguments/config
	if debugMode {
		restlerWrapperArgs = append(restlerWrapperArgs, "--debug")
	}

	restlerSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        jobName,
			Namespace:   namespace,
			Annotations: map[string]string{"cnfuzz/ignore": "true"},
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: "auth-script-map",
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: "auth-script",
									},
								},
							},
						},
					},
					Containers: []v1.Container{
						{
							Name:            containerName,
							Image:           containerImage,
							ImagePullPolicy: pullPolicy,
							Args:            restlerWrapperArgs,
							Env: []v1.EnvVar{
								{
									Name:  "RESTLER_TELEMETRY_OPTOUT",
									Value: telemetryOptOut,
								},
							},
							Resources: v1.ResourceRequirements{
								Limits: v1.ResourceList{
									v1.ResourceCPU:    cpuLimit,
									v1.ResourceMemory: memoryLimit,
								},
								Requests: v1.ResourceList{
									v1.ResourceCPU:    cpuRequest,
									v1.ResourceMemory: memoryRequest,
								},
							},
						},
					},
					ServiceAccountName: serviceAcc,
					RestartPolicy:      v1.RestartPolicyNever,
				},
			},
		},
	}
	return restlerSpec
}
