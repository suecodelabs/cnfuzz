package auth

import (
	"fmt"

	"github.com/suecodelabs/cnfuzz/src/discovery"
	"github.com/suecodelabs/cnfuzz/src/log"
)

type TokenSource interface {
	Token() (*Token, error)
}

func CreateTokenSource(schema discovery.SecuritySchema, clientId string, secret string) (TokenSource, error) {
	var createdTokenSource TokenSource
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
			createdTokenSource, err = CreateTokenFromOAuthFlow(flow.GrantType, clientId, secret, flow)
		}
		break
	default:
		// unkown security schema
		log.L().Infof("no tokensource available for %s auth scheme", schema.Key)
		return nil, fmt.Errorf("no tokensource available for %s auth scheme", schema.Key)
	}
	if err != nil {
		return nil, fmt.Errorf("error when creating a new token source: %w", err)
	}
	return createdTokenSource, nil
}

func CreateTokenSourceFromSchemas(schemas []discovery.SecuritySchema, clientId string, secret string) (TokenSource, error) {
	// Check if there are any security schemas
	if len(schemas) > 0 {
		// Just use the first one for now ...
		// If there are multiple schemas the situation does get weird, because the user can only pass one
		// username/secret combination, and there is currently no way to select a preferred auth scheme
		selScheme := 0
		for i, scheme := range schemas {
			if len(scheme.Flows) > 0 && scheme.Flows[0].GrantType == discovery.AuthorizationCode {
				selScheme = i
			}
		}
		selectedAuthScheme := schemas[selScheme]

		tokenSource, authErr := CreateTokenSource(selectedAuthScheme, clientId, secret)
		if authErr != nil {
			// Maybe if there are multiple security schemas, we could try a different one
			log.L().Errorf("error while creating an auth token: %+v", authErr)
		}
		return tokenSource, nil
	}
	return nil, fmt.Errorf("the API contains no security schemas")
}
