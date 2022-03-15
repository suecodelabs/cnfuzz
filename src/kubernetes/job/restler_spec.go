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

func createRestlerJob(fuzzConfig *config.FuzzConfig, targetPod *v1.Pod) *batchv1.Job {
	fullCommand := createRestlerCommand(targetPod, fuzzConfig)

	openApiVolumeName := "openapi-volume-" + fuzzConfig.KubernetesConfig.RestlerJobName
	initContainerUser := int64(0)
	volQuant := resource.MustParse("1Mi")

	restlerSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			// Maybe make this name unique?
			Name:        fuzzConfig.KubernetesConfig.RestlerJobName,
			Namespace:   fuzzConfig.KubernetesConfig.Namespace,
			Annotations: map[string]string{"cnfuzz/ignore": "true"},
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
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
							Name:  fuzzConfig.KubernetesConfig.RestlerInitJobName,
							Image: fuzzConfig.KubernetesConfig.RestlerInitImage,
							Args:  []string{fuzzConfig.ApiDescription.DiscoveryDoc.String(), "-s", "-S", "-o", "/openapi/doc.json"},
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
							Name:    fuzzConfig.KubernetesConfig.RestlerJobName,
							Image:   fuzzConfig.KubernetesConfig.RestlerImage,
							Command: []string{"/bin/sh", "-c"},
							Args:    []string{fullCommand},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      openApiVolumeName,
									MountPath: "/openapi",
								},
								{
									Name:      "auth-script-map",
									MountPath: "/scripts",
								},
							},
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
		},
	}
	return restlerSpec
}

func createRestlerCommand(targetPod *v1.Pod, fuzzConfig *config.FuzzConfig) string {
	targetIp := targetPod.Status.PodIP
	targetPort := fuzzConfig.ApiDescription.BaseUrl.Port()
	timeBudget := fuzzConfig.TimeBudget
	// Should we use SSL?
	isSsl := false
	if fuzzConfig.ApiDescription.BaseUrl.Scheme == "https" {
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
	token, authErr := createAuthToken(fuzzConfig)
	if authErr == nil && len(token) > 0 {
		// Use a high refresh interval because we have a static token (for now?)
		tokenCommand := fmt.Sprintf("--token_refresh_interval 999999 --token_refresh_command \"python3 /scripts/auth.py '%s' '%s'\"", targetPod.Name, token)
		fuzzCommand += " " + tokenCommand
	}
	fullCommand := fmt.Sprintf("%s && %s", compileCommand, fuzzCommand)

	return fullCommand
}

func createAuthToken(fuzzConfig *config.FuzzConfig) (token string, authErr error) {
	tokenSource, authErr := auth.CreateTokenSourceFromSchemas(fuzzConfig.ApiDescription.SecuritySchemes, fuzzConfig.ClientId, fuzzConfig.Secret)
	if authErr != nil {
		log.L().Errorf("error while building auth token provider: ", authErr)
		return "", authErr
	} else {
		tok, tokErr := tokenSource.Token()
		if tokErr != nil {
			log.L().Errorf("error while getting a new auth token: %+v", tokErr)
			return "", tokErr
		} else {
			token = fmt.Sprintf("%s: %s", "Authorization", tok.CreateAuthHeaderValue())
			return token, nil
		}
	}
}
