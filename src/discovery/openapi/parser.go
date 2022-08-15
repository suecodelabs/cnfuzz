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

package openapi

import (
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/suecodelabs/cnfuzz/src/logger"
	"net/url"
	"strconv"

	"github.com/getkin/kin-openapi/openapi2"
	conv "github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/suecodelabs/cnfuzz/src/discovery"
)

// UnMarshalOpenApiDoc unmarshal OpenAPI doc represented as a byte array
// returns a WebApiDescription object that represents the OpenAPI document
func UnMarshalOpenApiDoc(l logr.Logger, docFile []byte, uri *url.URL) (*discovery.WebApiDescription, error) {
	var doc *openapi3.T

	docIsVersion2 := false
	var err error
	version, err := getMajorDocVersion(l, docFile)
	if err != nil {
		return nil, fmt.Errorf("error while trying to read the used OpenAPI version in the retrieved doc: %w", err)
	} else {
		if version == 2 {
			docIsVersion2 = true
		} else if version == 3 {
			docIsVersion2 = false
		} else {
			return nil, fmt.Errorf("unknown OpenApi version")
		}
	}

	if docIsVersion2 {
		var doc2 openapi2.T
		if err := json.Unmarshal(docFile, &doc2); err != nil {
			return nil, fmt.Errorf("error while trying to unmarshal the OpenAPI doc: %w", err)
		}
		doc, err = conv.ToV3(&doc2)
		if err != nil {
			return nil, fmt.Errorf("error while trying to convert the OpenAPI v2 struct to a v3 struct: %w", err)
		}
	} else {
		doc, err = openapi3.NewLoader().LoadFromDataWithPath(docFile, uri)
		if err != nil {
			return nil, fmt.Errorf("error while loading OpenAPI doc from %s: %w", uri.String(), err)
		}
	}
	return parseOpenApiDoc(l, doc, uri, docIsVersion2)
}

// parseOpenApiDoc accepts a kin-openapi3 document and tries to convert it to a cnfuzz WebApiDescription object
func parseOpenApiDoc(l logr.Logger, doc *openapi3.T, uri *url.URL, docIsVersion2 bool) (*discovery.WebApiDescription, error) {
	// General info
	var desc discovery.WebApiDescription
	desc.Title = doc.Info.Title
	desc.Description = doc.Info.Description
	desc.Version = doc.Info.Version
	if uri == nil {
		return nil, fmt.Errorf("URL cant be empty or nil")
	}

	if docIsVersion2 {
		desc.DiscoverySource = discovery.OpenApiV2Source
	} else {
		desc.DiscoverySource = discovery.OpenApiV3Source
	}
	desc.DiscoveryDoc = *uri

	// Endpoints
	for strPath, pathObj := range doc.Paths {
		for method, operation := range pathObj.Operations() {
			endpoint := discovery.Endpoint{
				Method:      method,
				Path:        strPath,
				Description: operation.Description,
				Summary:     operation.Summary,
			}

			if operation.RequestBody != nil && operation.RequestBody.Value != nil {
				endpoint.Body = transformBody(l, operation.RequestBody.Value)
			}

			for _, paramObj := range operation.Parameters {
				if paramObj.Value == nil {
					l.V(logger.ImportantLevel).Info("a parameter in the OpenAPI doc is empty, this might be an invalid doc", "emptyParameter", paramObj.Ref)
				} else {
					resp := discovery.Parameter{
						Name:        paramObj.Value.Name,
						In:          paramObj.Value.In,
						Required:    paramObj.Value.Required,
						Description: paramObj.Value.Description,
						ParamType:   paramObj.Ref,
					}
					if paramObj.Value.Schema == nil || paramObj.Value.Schema.Value == nil {
						l.V(logger.DebugLevel).Info("no schema for parameter in OpenAPI doc", "parameter", paramObj.Ref)
					} else {
						resp.Schema = transformSchema(l, paramObj.Value.Schema.Ref, paramObj.Value.Schema.Value)
					}
					endpoint.Parameters = append(endpoint.Parameters, resp)
				}

			}

			for code, responseObj := range operation.Responses {
				codeInt, err := strconv.Atoi(code)
				if err != nil {
					l.V(logger.ImportantLevel).Info("Http status code in the OpenAPI doc is not a number", "statusCode", code)
					// Response object without a status code isn't very useful, so ignore it
					continue
				}
				resp := discovery.Response{
					Code: codeInt,
				}
				if responseObj.Value != nil {
					resp.Description = *responseObj.Value.Description
				}
				resp.Content = transformContent(l, responseObj.Value.Content)
				endpoint.Responses = append(endpoint.Responses, resp)
			}

			desc.Endpoints = append(desc.Endpoints, endpoint)
		}
	}

	// Security
	for key, scheme := range doc.Components.SecuritySchemes {
		schemeValue := scheme.Value
		newSchema := discovery.SecuritySchema{
			Key:              key,
			Type:             schemeValue.Type,
			Description:      schemeValue.Description,
			Name:             schemeValue.Name,
			In:               schemeValue.In,
			BearerFormat:     schemeValue.BearerFormat,
			OpenIdConnectUrl: schemeValue.OpenIdConnectUrl,
		}

		// Only oauth has flows
		if scheme.Value.Type == "oauth2" {
			// find flows
			if scheme.Value.Flows.Implicit != nil {
				foundFlow := scheme.Value.Flows.Implicit
				newFlow := transformOAuthFlow(discovery.Implicit, foundFlow)
				newSchema.Flows = append(newSchema.Flows, newFlow)
			}

			if scheme.Value.Flows.AuthorizationCode != nil {
				foundFlow := scheme.Value.Flows.AuthorizationCode
				newFlow := transformOAuthFlow(discovery.AuthorizationCode, foundFlow)
				newSchema.Flows = append(newSchema.Flows, newFlow)
			}

			if scheme.Value.Flows.Password != nil {
				foundFlow := scheme.Value.Flows.Password
				newFlow := transformOAuthFlow(discovery.Password, foundFlow)
				newSchema.Flows = append(newSchema.Flows, newFlow)
			}

			if scheme.Value.Flows.ClientCredentials != nil {
				foundFlow := scheme.Value.Flows.ClientCredentials
				newFlow := transformOAuthFlow(discovery.ClientCredentials, foundFlow)
				newSchema.Flows = append(newSchema.Flows, newFlow)
			}
		}

		desc.SecuritySchemes = append(desc.SecuritySchemes, newSchema)
	}

	return &desc, nil
}

// getMajorDocVersion tries to get the version of an OpenAPI doc
func getMajorDocVersion(l logr.Logger, doc []byte) (version int, err error) {
	var result map[string]any

	err = json.Unmarshal(doc, &result)
	if err != nil {
		return 0, err
	}

	// OpenAPI v2 spec uses "swagger" key for version
	swaggerVar := fmt.Sprint(result["swagger"])
	if len(swaggerVar) > 0 {
		// Probably version 1 or 2
		if swaggerVar[0] == '2' {
			return 2, nil
		} else if swaggerVar[0] == '1' {
			// We dont support this version
			return 0, fmt.Errorf("OpenAPI version 1 isn't supported")
		}
	}

	// OpenAPI v3 uses "openapi" key for version
	openapiVar := fmt.Sprint(result["openapi"])
	if len(openapiVar) > 0 {
		if openapiVar[0] == '3' {
			return 3, nil
		}
	}

	return 0, fmt.Errorf("version of the OpenAPI doc is unknown")
}

// transformOAuthFlow converts kin-openapi3 OAuthFlow to a cnfuzz OAuthFlow object
func transformOAuthFlow(grantType string, flow *openapi3.OAuthFlow) discovery.OAuthFlow {
	return discovery.OAuthFlow{
		GrantType:        grantType,
		AuthorizationURL: flow.AuthorizationURL,
		TokenURL:         flow.TokenURL,
		Scopes:           flow.Scopes,
		RefreshURL:       flow.RefreshURL,
	}
}

// transformBody converts kin-openapi3 RequestBody to a cnfuzz Body object
func transformBody(l logr.Logger, rBody *openapi3.RequestBody) discovery.Body {
	body := discovery.Body{
		Description: rBody.Description,
		Required:    rBody.Required,
	}
	if rBody.Content != nil {
		body.Content = transformContent(l, rBody.Content)
	}
	return body
}

// transformSchema converts kin-openapi3 Schema to a cnfuzz Schema object
func transformSchema(l logr.Logger, id string, schema *openapi3.Schema) discovery.Schema {
	schemaModel := discovery.Schema{
		Key:        id,
		Type:       schema.Type,
		Format:     schema.Format,
		Nullable:   schema.Nullable,
		AllowEmpty: schema.AllowEmptyValue,
		Example:    schema.Example,
	}

	for propId, schemaProp := range schema.Properties {
		if schemaProp == nil || schemaProp.Value == nil {
			l.V(logger.ImportantLevel).Info("schema property is nil or it's value is nil, the OpenAPI doc might be invalid", "schemaPropertyId", propId, "schemaProperty", schemaProp)
			continue
		}
		schemaModel.Properties = append(schemaModel.Properties, transformSchema(l, propId, schemaProp.Value))
	}

	return schemaModel
}

// transformContent converts kin-openapi3 Content to a cnfuzz Content object
func transformContent(l logr.Logger, contents openapi3.Content) []discovery.Content {
	if contents == nil {
		return nil
	}

	var responses []discovery.Content
	for contentType, schemaRef := range contents {
		content := discovery.Content{
			ContentType: contentType,
		}

		content.Schema = transformSchema(l, schemaRef.Schema.Ref, schemaRef.Schema.Value)

		responses = append(responses, content)
	}
	return responses
}
