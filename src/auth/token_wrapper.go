package auth

import (
	"fmt"

	"github.com/suecodelabs/cnfuzz/src/auth/source"
	"github.com/suecodelabs/cnfuzz/src/discovery"
	"github.com/suecodelabs/cnfuzz/src/log"
)

type TokenWrapper struct {
	Schema   discovery.SecuritySchema
	ClientId string
	Secret   string

	TokenSource source.TokenSource
}

func CreateTokenWrapper(schema discovery.SecuritySchema, clientId string, secret string) (*TokenWrapper, error) {
	tokenWrapper := TokenWrapper{
		Schema:   schema,
		ClientId: clientId,
		Secret:   secret,
	}
	err := tokenWrapper.CreateTokenSource()
	if err != nil {
		return nil, err
	}
	return &tokenWrapper, nil
}

func CreateTokenWrapperFromSchema(schemes []discovery.SecuritySchema, clientId string, secret string) (*TokenWrapper, error) {
	// Auth
	var tokenWrapper *TokenWrapper
	// Check if there are any security schemes
	if len(schemes) > 0 {
		// Just use the first one for now ...
		// If there are multiple schemes the situation does get weird, because the user can only pass one
		// username/secret combination, and there is currently no way to select a preferred auth scheme
		selScheme := 0
		for i, scheme := range schemes {
			if len(scheme.Flows) > 0 && scheme.Flows[0].GrantType == discovery.AuthorizationCode {
				selScheme = i
			}
		}
		selectedAuthScheme := schemes[selScheme]

		var authErr error
		tokenWrapper, authErr = CreateTokenWrapper(selectedAuthScheme, clientId, secret)
		if authErr != nil {
			// Maybe if there are multiple security schemes, we could try a different one
			log.L().Errorf("error while creating an auth token: %+v", authErr)
		}
		return tokenWrapper, nil
	}
	return nil, fmt.Errorf("the API contains no security schemes")
}

func (tokenWrapper *TokenWrapper) CreateToken(scheme discovery.SecuritySchema, clientId string, secret string) string {
	tok, err := tokenWrapper.TokenSource.Token()
	if err == nil {
		log.L().Errorf("error while creating a new auth access token: %+v", err)
	}
	return tok.AccessToken
}

func (tokenWrapper *TokenWrapper) CreateTokenSource() error {
	var createdTokenSource source.TokenSource
	var err error
	switch tokenWrapper.Schema.Type {
	case discovery.BasicSecSchemaType:
		createdTokenSource, err = source.BasicAuthTokenSource(tokenWrapper.ClientId, tokenWrapper.Secret)
		break
	case discovery.ApiKeySecSchemaType:
		createdTokenSource, err = source.ApiKeyTokenSource(tokenWrapper.Secret)
		break
	case discovery.OAuth2SecSchemaType:
		for _, flow := range tokenWrapper.Schema.Flows {
			createdTokenSource, err = source.CreateTokenFromOAuthFlow(flow.GrantType, tokenWrapper.ClientId, tokenWrapper.Secret, flow)
		}
		break
	default:
		// unkown security schema
		log.L().Infof("no tokensource available for %s auth scheme", tokenWrapper.Schema.Key)
		return fmt.Errorf("no tokensource available for %s auth scheme", tokenWrapper.Schema.Key)
	}
	if err != nil {
		return fmt.Errorf("error when creating a new token source: %w", err)
	}
	tokenWrapper.TokenSource = createdTokenSource
	return nil
}
