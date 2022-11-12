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

package restlerwrapper

import (
	"github.com/suecodelabs/cnfuzz/src/pkg/discovery"
	"github.com/suecodelabs/cnfuzz/src/pkg/discovery/openapi"
	"github.com/suecodelabs/cnfuzz/src/pkg/logger"
	"github.com/suecodelabs/cnfuzz/src/pkg/restlerwrapper/auth"
	"os"
)

type SomeInfo struct {
	ApiDesc     *discovery.WebApiDescription
	TokenSource auth.ITokenSource
}

func CollectInfoFromAddr(l logger.Logger, ip string, ports []int32, oaLocs []string, dryRun bool) SomeInfo {
	apiDoc, err := openapi.TryGetOpenApiDoc(l, ip, ports, oaLocs)
	if err != nil {
		l.FatalError(err, "error while retrieving OpenAPI document")
	}

	apiDesc, err := openapi.ParseOpenApiDoc(l, apiDoc)
	if err != nil {
		l.FatalError(err, "error while unmarshalling OpenAPI doc request body")
	}

	b, err := apiDoc.DocFile.MarshalJSON()
	if err != nil {
		l.FatalError(err, "failed to marshal OpenApi doc to bytes")
	} else {
		if !dryRun {
			err := os.Mkdir("/openapi", os.FileMode(0755))
			if err != nil {
				l.FatalError(err, "failed to create 'openapi' dir to write OpenApi doc into")
			}
			err = os.WriteFile("/openapi/doc.json", b, os.FileMode(0644))
			if err != nil {
				l.FatalError(err, "failed to write OpenApi doc to fs")
			}
		}
	}

	// Tokensource can be nil !!! this means the API doesn't have any security (specified in the discovery doc ...)
	tokenSource, authErr := auth.CreateTokenSourceFromSchemas(l, apiDesc.SecuritySchemes, "username", "secret") // TODO cnf.AuthConfig.Username, cnf.AuthConfig.Secret)
	if authErr != nil {
		l.FatalError(authErr, "error while building auth token source")
	}

	return SomeInfo{
		TokenSource: tokenSource,
		ApiDesc:     apiDesc,
	}
}
