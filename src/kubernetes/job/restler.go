package job

import (
	"fmt"

	"github.com/suecodelabs/cnfuzz/src/auth"
	"github.com/suecodelabs/cnfuzz/src/config"
	"github.com/suecodelabs/cnfuzz/src/log"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// createRestlerJob creates a Kubernetes Job for the RESTler fuzzer
// this includes an init container that gets the OpenAPI doc from the target API with curl and volumes for transferring the information
// it uses values from the FuzzConfig to configure the fuzz command that runs inside the RESTler container
func createRestlerJob(cnf *config.FuzzerConfig, tokenSource auth.ITokenSource) *batchv1.Job {
	reportDir := "/reportdir"
	// FIXME don't hardcode the endpoint url

	restlerCommand := createRestlerCommand(cnf, tokenSource, reportDir)
	awsCliCommand := createAwsCliCommand(cnf.ProcessResultConf, reportDir)

	reportVolumeName := "result-volume-" + cnf.JobName
	openApiVolumeName := "openapi-volume-" + cnf.JobName
	initContainerUser := int64(0)
	volQuant := resource.MustParse("1Mi")

	restlerSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			// Maybe make this name unique?
			Name:        cnf.JobName,
			Namespace:   cnf.Namespace,
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
							Name:  cnf.InitJobName,
							Image: cnf.InitImage,
							Args:  []string{cnf.DiscoveryDocLocation, "-s", "-S", "-o", "/openapi/doc.json"},
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
							Name:    cnf.JobName,
							Image:   cnf.Image,
							Command: []string{"/bin/sh", "-c"},
							Args:    []string{restlerCommand},
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
							Name:    cnf.ProcessResultConf.ContainerName,
							Image:   cnf.ProcessResultConf.Image,
							Command: []string{"/bin/sh", "-c"},
							Args:    []string{awsCliCommand},
							Env: []v1.EnvVar{
								{
									Name:  "AWS_ACCESS_KEY_ID",
									Value: cnf.ProcessResultConf.AccessKey,
								},
								{
									Name:  "AWS_SECRET_ACCESS_KEY",
									Value: cnf.ProcessResultConf.SecretKey,
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
func createRestlerCommand(cnf *config.FuzzerConfig, tokenSource auth.ITokenSource, reportVol string) string {
	targetIp := cnf.Target.IP
	targetPort := cnf.Target.Port
	timeBudget := cnf.TimeBudget
	// Should we use SSL?
	isSsl := false
	if cnf.Target.Scheme == "https" {
		log.L().Debug("using SSL")
		isSsl = true
	} else {
		log.L().Debug("not using SSL")
	}

	log.L().Debugf("using target_ip %s and target_port %s", targetIp, targetPort)
	compileCommand := fmt.Sprintf("dotnet /RESTler/restler/Restler.dll compile --api_spec /openapi/doc.json")
	// Please, UNIX philosophy people.
	fuzzCommand := fmt.Sprintf("dotnet /RESTler/restler/Restler.dll fuzz --grammar_file /Compile/grammar.py --dictionary_file /Compile/dict.json --target_ip %s --target_port %s --time_budget %s", targetIp, targetPort, timeBudget)
	if !isSsl {
		fuzzCommand = fmt.Sprintf("%s --no_ssl", fuzzCommand)
	}
	// create a new auth token using the tokensource
	tok, tokErr := tokenSource.Token()
	if tokErr != nil {
		log.L().Errorf("error while getting a new auth token: %+v", tokErr)
	} else {
		token := fmt.Sprintf("%s: %s", "Authorization", tok.CreateAuthHeaderValue())
		if tokErr == nil && len(token) > 0 {
			// Use a high refresh interval because we have a static token (for now?)
			tokenCommand := fmt.Sprintf("--token_refresh_interval 999999 --token_refresh_command \"python3 /scripts/auth.py '%s' '%s'\"", cnf.Target.PodName, token)
			fuzzCommand += " " + tokenCommand
		}
	}
	// FIXME I think the fuzz directory might be called fuzlean when fuzzing in lean mode but haven't checked yet
	// FIXME move this towards PreStop lifecycle hook of pod
	copyCommand := fmt.Sprintf("mv /Fuzz/* %s", reportVol)

	fullCommand := fmt.Sprintf("%s && %s && %s", compileCommand, fuzzCommand, copyCommand)

	return fullCommand
}

func createAwsCliCommand(cnf config.S3Config, reportMountDir string) string {
	baseAwsCmd := "aws s3"
	if len(cnf.EndpointUrl) > 0 {
		baseAwsCmd = fmt.Sprintf("aws --endpoint-url %s s3", cnf.EndpointUrl)
	}

	waitCommand := fmt.Sprintf("until (( $(ls -1q %s | wc -l) > 1 )); do sleep 5; done;", reportMountDir)
	createBucketCommand := fmt.Sprintf("%s mb %s", baseAwsCmd, cnf.ReportBucket)
	copyCommand := fmt.Sprintf("%s cp --recursive %s %s", baseAwsCmd, reportMountDir, cnf.ReportBucket)

	fullCommand := fmt.Sprintf("%s %s; %s", waitCommand, createBucketCommand, copyCommand)

	return fullCommand
}
