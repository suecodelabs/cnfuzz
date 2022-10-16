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

package auth

import (
	"fmt"
	"github.com/suecodelabs/cnfuzz/src/logger"

	"github.com/suecodelabs/cnfuzz/src/discovery"
)

// ITokenSource interface for creating new auth tokens
type ITokenSource interface {
	Token() (*Token, error)
}

// CreateTokenSource creates a new ITokenSource
// Uses the schema type (BasicSecSchemaType) to create a ITokenSource for the proper auth source
func CreateTokenSource(l logger.Logger, schema discovery.SecuritySchema, clientId string, secret string) (ITokenSource, error) {
	var createdTokenSource ITokenSource
	var err error
	switch schema.Type {
	case discovery.BasicSecSchemaType:
		createdTokenSource, err = BasicAuthTokenSource(clientId, secret)
		break
	case discovery.ApiKeySecSchemaType:
		createdTokenSource, err = ApiKeyTokenSource(secret)
		break
	case discovery.OAuth2SecSchemaType:
		for _, flow := range schema.Flows {
			createdTokenSource, err = CreateTokenFromOAuthFlow(l, flow.GrantType, clientId, secret, flow)
		}
		break
	default:
		// unkown security schema
		l.V(logger.ImportantLevel).Info("no token source available for auth scheme", "authScheme", schema.Key)
		return nil, fmt.Errorf("no tokensource available for %s auth scheme", schema.Key)
	}
	if err != nil {
		return nil, fmt.Errorf("error when creating a new token source: %w", err)
	}
	return createdTokenSource, nil
}

// CreateTokenSourceFromSchemas creates a new ITokenSource from the first schema in the slice
func CreateTokenSourceFromSchemas(l logger.Logger, schemas []discovery.SecuritySchema, clientId string, secret string) (ITokenSource, error) {
	// Check if there are any security schemas
	// This function could be improved by having a smarter algorithm for picking a schema
	if len(schemas) > 0 {
		// Just use the first one for now ...
		// If there are multiple schemas the situation does get weird, because the user can only pass one
		// username/secret combination, and there is currently no way to select a preferred auth scheme
		selScheme := 0
		/* for i, scheme := range schemas {
			if len(scheme.Flows) > 0 && scheme.Flows[0].GrantType == discovery.AuthorizationCode {
				selScheme = i
			}
		} */
		selectedAuthScheme := schemas[selScheme]

		tokenSource, authErr := CreateTokenSource(l, selectedAuthScheme, clientId, secret)
		if authErr != nil {
			// Maybe if there are multiple security schemas, we could try a different one
			l.V(logger.ImportantLevel).Error(authErr, "error while creating an auth token")
		}
		return tokenSource, nil
	}
	return nil, nil
}
