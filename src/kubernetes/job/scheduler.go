package job

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/suecodelabs/cnfuzz/src/cmd"
	"github.com/suecodelabs/cnfuzz/src/config"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// createSchedulerJob creates a Kubernetes Job for the cnfuzz instance that gets the OpenAPI spec and starts the RESTler fuzzer job
func createSchedulerJob(config *config.SchedulerConfig) *batchv1.Job {
	var backOffLimit int32 = 1
	var terminateAfter int64 = 240
	var TTLAfterFinish int32 = 120

	// Args for the job
	args := buildSchedulerArgs(config)

	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        config.JobName,
			Namespace:   config.Namespace,
			Annotations: map[string]string{"cnfuzz/ignore": "true"},
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  config.JobName,
							Image: config.Image,
							Args:  args,
							// Command: strings.Split(*cmd, " "),
						},
					},
					ServiceAccountName: "cnfuzz-job",
					RestartPolicy:      v1.RestartPolicyNever,
				},
			},
			ActiveDeadlineSeconds:   &terminateAfter,
			BackoffLimit:            &backOffLimit,
			TTLSecondsAfterFinished: &TTLAfterFinish,
		},
	}
	return jobSpec
}

// buildSchedulerArgs creates a string array with flags/arguments for the cnfuzz job
func buildSchedulerArgs(config *config.SchedulerConfig) []string {
	// Args for the job
	podNameArg := fmt.Sprintf("--%s", cmd.TargetPodName)
	podNameVal := config.TargetPodName
	podNamespaceArg := fmt.Sprintf("--%s", cmd.TargetPodNamespace)
	podNamespaceVal := config.TargetPodNamespace

	args := []string{podNameArg, podNameVal, podNamespaceArg, podNamespaceVal}

	if viper.GetBool(cmd.IsDebug) {
		args = append(args, "--debug")
	}

	stringFlagsToPassDown := []string{cmd.AuthUsername, cmd.AuthSecretFlag, cmd.HomeNamespaceFlag, cmd.RestlerTimeBudget}
	for _, arg := range stringFlagsToPassDown {
		setValue := viper.GetString(arg)
		if len(setValue) == 0 {
			continue
		}
		args = append(args, fmt.Sprintf("--%s", arg))
		args = append(args, setValue)
	}

	return args
}
