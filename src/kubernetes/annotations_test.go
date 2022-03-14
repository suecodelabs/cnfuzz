package kubernetes

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/suecodelabs/cnfuzz/src/cmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"testing"
)

func TestSetConfigRegister(t *testing.T) {
	uname := "testusername"
	secret := "verysecret"
	testAnnos := Annotations{
		Username: uname,
		Secret:   secret,
	}

	testAnnos.SetConfigRegister()

	assert.Equal(t, uname, viper.GetString(cmd.AuthUsername))
	assert.Equal(t, secret, viper.GetString(cmd.AuthSecretFlag))
}

func TestGetAnnotations(t *testing.T) {
	ignoreMeVal := true
	fuzzMeVal := false
	oaDocVal := "/swagger/swagger.json"
	secretVal := "verysecret"
	unameVal := "me"
	testMeta := &metav1.ObjectMeta{
		Annotations: map[string]string{
			fmt.Sprintf("%s/%s", AnnotationPrefix, IgnoreMeAnno):   strconv.FormatBool(ignoreMeVal),
			fmt.Sprintf("%s/%s", AnnotationPrefix, FuzzMeAnno):     strconv.FormatBool(fuzzMeVal),
			fmt.Sprintf("%s/%s", AnnotationPrefix, OpenApiDocAnno): oaDocVal,
			fmt.Sprintf("%s/%s", AnnotationPrefix, SecretAnno):     secretVal,
			fmt.Sprintf("%s/%s", AnnotationPrefix, UsernameAnno):   unameVal,
		},
	}
	result := GetAnnotations(testMeta)
	assert.Equal(t, ignoreMeVal, result.IgnoreMe)
	assert.Equal(t, fuzzMeVal, result.FuzzMe)
	assert.Equal(t, oaDocVal, result.OpenApiDocLocation)
	assert.Equal(t, secretVal, result.Secret)
	assert.Equal(t, unameVal, result.Username)
}

func TestGetAnnotationsFromMeta(t *testing.T) {
	testAnno := "test"
	testValue := "myvalue"
	testMeta := &metav1.ObjectMeta{
		Annotations: map[string]string{
			fmt.Sprintf("%s/%s", AnnotationPrefix, testAnno): testValue,
		},
	}
	result := getAnnotationFromMeta(testMeta, testAnno)
	assert.Equal(t, testValue, result)
}
