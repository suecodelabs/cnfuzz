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
	"github.com/go-logr/logr"
	"github.com/suecodelabs/cnfuzz/src/discovery"
	"github.com/suecodelabs/cnfuzz/src/logger"
	"time"

	"github.com/suecodelabs/cnfuzz/src/auth"
	"github.com/suecodelabs/cnfuzz/src/config"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type restlerFuzzJob struct {
	JobName              string
	Namespace            string
	ContainerName        string
	ContainerImage       string
	InitContainerName    string
	InitContainerImage   string
	DiscoveryDocLocation string
	TelemetryOptOut      string
	// TimeBudget in hours https://github.com/microsoft/restler-fuzzer/blob/1df2d3b8177408dd04666926867e10d505249281/docs/user-guide/SettingsFile.md#time_budget-float-default-30-days
	TimeBudget    string
	CpuLimit      int64
	MemoryLimit   int64
	CpuRequest    int64
	MemoryRequest int64
	Target        *restlerFuzzTarget
	S3Config      *config.S3Config
}

type restlerFuzzTarget struct {
	podName   string
	namespace string
	// IP of the target https://github.com/microsoft/restler-fuzzer/blob/1df2d3b8177408dd04666926867e10d505249281/docs/user-guide/SettingsFile.md#target_ip-str-default-none
	ip string
	// Port of the target https://github.com/microsoft/restler-fuzzer/blob/1df2d3b8177408dd04666926867e10d505249281/docs/user-guide/SettingsFile.md#target_port-int-default-none
	port string
	// Http scheme (http, https)
	scheme string
}

func CreateRestlerJob(targetPod *v1.Pod, cnf *config.CnFuzzConfig, discoveryDoc *discovery.WebApiDescription) *restlerFuzzJob {
	j := restlerFuzzJob{
		JobName:              "cnfuzz-job-" + targetPod.Name,
		Namespace:            targetPod.Namespace, // TODO do we ever not use the target pod namespace?
		ContainerName:        "cnfuzz-job-" + targetPod.Name + "-restler",
		ContainerImage:       cnf.RestlerConfig.Image,
		InitContainerName:    "cnfuzz-job-" + targetPod.Name + "-restler-init",
		InitContainerImage:   cnf.RestlerConfig.InitImage,
		DiscoveryDocLocation: discoveryDoc.DiscoveryDoc.String(),
		TelemetryOptOut:      cnf.RestlerConfig.TelemetryOptOut,
		TimeBudget:           cnf.RestlerConfig.TimeBudget,
		CpuLimit:             cnf.RestlerConfig.CpuLimit,
		MemoryLimit:          cnf.RestlerConfig.MemoryLimit,
		Target: &restlerFuzzTarget{
			podName:   targetPod.Name,
			namespace: targetPod.Namespace,
			ip:        targetPod.Status.PodIP,
			port:      discoveryDoc.DiscoveryDoc.Port(),
			scheme:    discoveryDoc.DiscoveryDoc.Scheme,
		},
		S3Config: cnf.S3Config,
	}

	return &j
}

// CreateRestlerJob creates a Kubernetes Job for the RESTler fuzzer
// this includes an init container that gets the OpenAPI doc from the target API with curl and volumes for transferring the information
// it uses values from the FuzzConfig to configure the fuzz command that runs inside the RESTler container
// returned job hasn't started yet
func (job *restlerFuzzJob) CreateRestlerJob(l logr.Logger, tokenSource auth.ITokenSource) *batchv1.Job {
	reportDir := "/reportdir"
	// File that when created triggers the S3 CLI container to copy the reports to the S3 bucket
	moveTriggerFile := reportDir + "/move_alert"

	fullCommand := job.createRestlerCommand(l, tokenSource, reportDir, moveTriggerFile)

	timeStamp := time.Now().Format("20060102150405")
	targetReportDir := fmt.Sprintf("%s/%s/%s", job.S3Config.ReportBucket, job.Target.podName, timeStamp)
	awsCliCommand := job.createAwsCliCommand(l, reportDir, targetReportDir, moveTriggerFile)

	reportVolumeName := "result-volume-" + job.JobName
	openApiVolumeName := "openapi-volume-" + job.JobName
	initContainerUser := int64(0)
	volQuant := resource.MustParse("1Mi")

	restlerSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        job.JobName,
			Namespace:   job.Namespace,
			Annotations: map[string]string{"cnfuzz/ignore": "true"},
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: reportVolumeName,
							VolumeSource: v1.VolumeSource{
								EmptyDir: &v1.EmptyDirVolumeSource{
									SizeLimit: &volQuant,
								},
							},
						},
						{
							Name: openApiVolumeName,
							VolumeSource: v1.VolumeSource{
								EmptyDir: &v1.EmptyDirVolumeSource{
									SizeLimit: &volQuant,
								},
							},
						},
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
					InitContainers: []v1.Container{
						{
							Name:  job.InitContainerName,
							Image: job.InitContainerImage,
							Args:  []string{job.DiscoveryDocLocation, "-s", "-S", "-o", "/openapi/doc.json"},
							SecurityContext: &v1.SecurityContext{
								RunAsUser: &initContainerUser,
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      openApiVolumeName,
									MountPath: "/openapi",
								},
							},
						},
					},
					Containers: []v1.Container{
						{
							Name:    job.ContainerName,
							Image:   job.ContainerImage,
							Command: []string{"/bin/sh", "-c"},
							Args:    []string{fullCommand},
							Env: []v1.EnvVar{
								{
									Name:  "RESTLER_TELEMETRY_OPTOUT",
									Value: job.TelemetryOptOut,
								},
							},
							Resources: v1.ResourceRequirements{
								Limits: v1.ResourceList{
									// CPU, in cores. (500m = .5 cores)
									v1.ResourceCPU: *resource.NewMilliQuantity(job.CpuLimit, resource.DecimalSI),
									// Memory, in bytes. (500Mi = 500MiB = 500 * 1024 * 1024)
									v1.ResourceMemory: *resource.NewQuantity(job.MemoryLimit*(1024*1024), resource.DecimalSI),
								},
								Requests: v1.ResourceList{
									// CPU, in cores. (500m = .5 cores)
									v1.ResourceCPU: *resource.NewMilliQuantity(job.CpuRequest, resource.DecimalSI),
									// Memory, in bytes. (500Mi = 500MiB = 500 * 1024 * 1024)
									v1.ResourceMemory: *resource.NewQuantity(job.MemoryRequest*(1024*1024), resource.DecimalSI),
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      openApiVolumeName,
									MountPath: "/openapi",
								},
								{
									Name:      reportVolumeName,
									MountPath: reportDir,
								},
								{
									Name:      "auth-script-map",
									MountPath: "/scripts",
								},
							},
						},
						{
							Name:    job.S3Config.Sidecar.ContainerName,
							Image:   job.S3Config.Sidecar.Image,
							Command: []string{"/bin/sh", "-c"},
							Args:    []string{awsCliCommand},
							Env: []v1.EnvVar{
								{
									Name:  "AWS_ACCESS_KEY_ID",
									Value: job.S3Config.AccessKey,
								},
								{
									Name:  "AWS_SECRET_ACCESS_KEY",
									Value: job.S3Config.SecretKey,
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      reportVolumeName,
									MountPath: reportDir, // this doesn't have to be the same dir as restler
								},
							},
							ImagePullPolicy: v1.PullIfNotPresent,
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
		},
	}
	return restlerSpec
}

// createRestlerCommand creates command string that can be run inside the RESTler container
// the command string consists of a compile command that analyzes the OpenAPI spec and generates a fuzzing grammar
// and the fuzz command itself
func (job *restlerFuzzJob) createRestlerCommand(l logr.Logger, tokenSource auth.ITokenSource, reportVol string, moveTriggerFile string) string {
	targetIp := job.Target.ip
	targetPort := job.Target.port
	timeBudget := job.TimeBudget
	// Should we use SSL?
	isSsl := false
	if job.Target.scheme == "https" {
		l.V(logger.DebugLevel).Info("using SSL in Restler")
		isSsl = true
	} else {
		l.V(logger.DebugLevel).Info("not using SSL in Restler")
	}

	l.V(logger.DebugLevel).Info(fmt.Sprintf("using %s:%s for restler", targetIp, targetPort), "targetIp", targetIp, "targetPort", targetPort)
	compileCommand := fmt.Sprintf("dotnet /RESTler/restler/Restler.dll compile --api_spec /openapi/doc.json")
	// Please, UNIX philosophy people.
	fuzzCommand := fmt.Sprintf("dotnet /RESTler/restler/Restler.dll fuzz --grammar_file /Compile/grammar.py --dictionary_file /Compile/dict.json --target_ip %s --target_port %s --time_budget %s", targetIp, targetPort, timeBudget)
	if !isSsl {
		fuzzCommand = fmt.Sprintf("%s --no_ssl", fuzzCommand)
	}

	if tokenSource != nil {
		// create a new auth token using the tokensource
		tok, tokErr := tokenSource.Token()
		if tokErr != nil {
			l.V(logger.ImportantLevel).Error(tokErr, "error while getting a new auth token")
		} else {
			token := fmt.Sprintf("%s: %s", "Authorization", tok.CreateAuthHeaderValue(l))
			if tokErr == nil && len(token) > 0 {
				// Use a high refresh interval because we have a static token (for now?)
				tokenCommand := fmt.Sprintf("--token_refresh_interval 999999 --token_refresh_command \"python3 /scripts/auth.py '%s' '%s'\"", job.Target.podName, token)
				fuzzCommand += " " + tokenCommand
			}
		}
	}

	// FIXME I think the fuzz directory might be called fuzzlean when fuzzing in lean mode but haven't checked yet
	// FIXME move this towards PreStop lifecycle hook of pod
	copyCommand := fmt.Sprintf("mv /Fuzz/* %s", reportVol)
	triggerCommand := fmt.Sprintf("touch %s", moveTriggerFile)

	fullCommand := fmt.Sprintf("%s && %s && %s && %s", compileCommand, fuzzCommand, copyCommand, triggerCommand)

	return fullCommand
}

func (job *restlerFuzzJob) createAwsCliCommand(l logr.Logger, reportMountDir string, targetReportDir string, triggerFile string) string {
	baseAwsCmd := "aws s3"
	if len(job.S3Config.EndpointUrl) > 0 {
		baseAwsCmd = fmt.Sprintf("aws --endpoint-url %s s3", job.S3Config.EndpointUrl)
	}

	waitCommand := fmt.Sprintf("until [ -f %s ]; do sleep 5; done;", triggerFile)
	copyCommand := fmt.Sprintf("%s cp --recursive %s %s", baseAwsCmd, reportMountDir, targetReportDir)

	fullCommand := fmt.Sprintf("%s %s", waitCommand, copyCommand)

	return fullCommand
}
