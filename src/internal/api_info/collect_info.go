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

package api_info

import (
	"context"
	"fmt"
	"github.com/suecodelabs/cnfuzz/src/pkg/auth"
	"github.com/suecodelabs/cnfuzz/src/pkg/config"
	"github.com/suecodelabs/cnfuzz/src/pkg/discovery"
	"github.com/suecodelabs/cnfuzz/src/pkg/discovery/openapi"
	"github.com/suecodelabs/cnfuzz/src/pkg/k8s"
	"github.com/suecodelabs/cnfuzz/src/pkg/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

// TargetInfo info about a target pod for fuzzing.
type TargetInfo struct {
	TargetAddr     string
	Annos          k8s.Annotations
	ApiDesc        *discovery.WebApiDescription
	UnparsedApiDoc openapi.UnParsedOpenApiDoc
	TokenSource    auth.ITokenSource
}

// CollectInfo collects info from target pod.
// returns TargetInfo.
func CollectInfo(l logger.Logger, targetPodName, targetNamespace, dDocIp string, dDocLoc string, ports []int32) TargetInfo {
	l.V(logger.DebugLevel).Info("getting pod info")
	pod := GetPod(l, targetPodName, targetNamespace, config.RunCnf.LocalK8sConfig)
	targetAddr := fmt.Sprintf("%s.%s.pod", strings.ReplaceAll(pod.Status.PodIP, ".", "-"), pod.Namespace)
	annos := k8s.GetAnnotations(&pod.ObjectMeta)

	l.V(logger.DebugLevel).Info("getting OpenApi document")
	var oaLocs []string
	if len(dDocLoc) > 0 {
		oaLocs = append(oaLocs, dDocLoc)
	} else if len(annos.OpenApiDocLocation) > 0 {
		oaLocs = append(oaLocs, annos.OpenApiDocLocation)
	} else {
		oaLocs = openapi.GetCommonOpenApiLocations()
	}

	oAAddr := targetAddr
	if len(dDocIp) > 0 {
		oAAddr = dDocIp
	}
	apiDoc, apiDesc := GetOpenApiDoc(l, oAAddr, ports, oaLocs)
	l.V(logger.DebugLevel).Info("found OpenApi document")

	l.V(logger.DebugLevel).Info("creating auth token source from pod annotations and OpenApi document")
	tokenSource := CreateTokenSource(l, apiDesc, annos.Username, annos.Secret)

	return TargetInfo{
		TargetAddr:     targetAddr,
		Annos:          annos,
		ApiDesc:        apiDesc,
		UnparsedApiDoc: apiDoc,
		TokenSource:    tokenSource,
	}
}

// GetPod get a Pod struct with a pod name and namespace.
func GetPod(l logger.Logger, targetPodName, targetNamespace string, useLocalConfig bool) *corev1.Pod {
	// parse the passed arguments
	var podName string
	if len(targetPodName) > 0 {
		podName = targetPodName
	} else {
		l.Fatal("no target given, pod name is empty")
	}

	client := k8s.CreateClientset(l, !useLocalConfig)
	pod, err := client.CoreV1().Pods(targetNamespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		l.FatalError(err, "failed to find target pod")
	}
	return pod
}

// GetOpenApiDoc get the OpenApi doc in discovery.WebApiDescription struct from a target host.
func GetOpenApiDoc(l logger.Logger, host string, ports []int32, oaLocs []string) (openapi.UnParsedOpenApiDoc, *discovery.WebApiDescription) {
	apiDoc, err := openapi.TryGetOpenApiDoc(l, host, ports, oaLocs)
	if err != nil {
		l.FatalError(err, "error while retrieving OpenAPI document")
	}
	apiDesc, err := openapi.ParseOpenApiDoc(l, apiDoc)
	if err != nil {
		l.FatalError(err, "error while unmarshalling OpenAPI doc request body")
	}

	return apiDoc, apiDesc
}

// CreateTokenSource creates a auth.ITokenSource from a discovery.WebApiDescription, username and secret.
// auth.ITokenSource that this function returns can be nil. This happens when the API doesn't have any security info specified.
func CreateTokenSource(l logger.Logger, apiDesc *discovery.WebApiDescription, username, secret string) auth.ITokenSource {
	// Tokensource can be nil !!! this means the API doesn't have any security (specified in the discovery doc ...)
	tokenSource, authErr := auth.CreateTokenSourceFromSchemas(l, apiDesc.SecuritySchemes, username, secret) // TODO cnf.AuthConfig.Username, cnf.AuthConfig.Secret)
	if authErr != nil {
		l.FatalError(authErr, "error while building auth token source")
	}

	return tokenSource
}
