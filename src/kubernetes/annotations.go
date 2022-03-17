package kubernetes

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
	"github.com/suecodelabs/cnfuzz/src/cmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	AnnotationPrefix = "cnfuzz"
	IgnoreMeAnno     = "ignore"
	FuzzMeAnno       = "enable"
	OpenApiDocAnno   = "open-api-doc"
	SecretAnno       = "secret"
	UsernameAnno     = "username"
)

// Annotations annotation values for annotations to be used inside Kubernetes configurations
type Annotations struct {
	IgnoreMe           bool
	FuzzMe             bool
	OpenApiDocLocation string
	Secret             string
	Username           string
}

// SetConfigRegister looks in the annotations object for empty values and tries to fill them with values from the config register
func (annos Annotations) SetConfigRegister() {
	if len(annos.Username) > 0 {
		viper.Set(cmd.AuthUsername, annos.Username)
	}

	if len(annos.Secret) > 0 {
		viper.Set(cmd.AuthSecretFlag, annos.Secret)
	}
}

// GetAnnotations gather annotations inside the metadata of a Kubernetes object
func GetAnnotations(objectMeta *metav1.ObjectMeta) Annotations {
	strIgnoreMe := getAnnotationFromMeta(objectMeta, IgnoreMeAnno)
	strFuzzMe := getAnnotationFromMeta(objectMeta, FuzzMeAnno)
	oaDocLoc := getAnnotationFromMeta(objectMeta, OpenApiDocAnno)
	secret := getAnnotationFromMeta(objectMeta, SecretAnno)
	username := getAnnotationFromMeta(objectMeta, UsernameAnno)

	ignoreMe, err := strconv.ParseBool(strIgnoreMe)
	if err != nil {
		// The value is invalid
		// just ignore the annotation
		ignoreMe = false
	}

	fuzzMe, err := strconv.ParseBool(strFuzzMe)

	return Annotations{
		IgnoreMe:           ignoreMe,
		FuzzMe:             fuzzMe,
		OpenApiDocLocation: oaDocLoc,
		Secret:             secret,
		Username:           username,
	}

}

// getAnnotationFromMeta get a single annotation value from Kubernetes object meta
func getAnnotationFromMeta(objectMeta *metav1.ObjectMeta, annotationName string) string {
	return objectMeta.Annotations[fmt.Sprintf("%s/%s", AnnotationPrefix, annotationName)]
}
