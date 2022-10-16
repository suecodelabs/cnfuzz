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
	"context"
	"errors"
	"fmt"
	"github.com/suecodelabs/cnfuzz/src/logger"

	"github.com/suecodelabs/cnfuzz/src/discovery"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// CreateTokenFromOAuthFlow creates a ITokenSource for a OAuthFlow
func CreateTokenFromOAuthFlow(l logger.Logger, grantType string, clientId string, secret string, flow discovery.OAuthFlow) (ITokenSource, error) {
	ctx := context.TODO()

	// Just getting all the possible scopes
	// This might cause problems, should allow for a scope collection to be passed with a flag
	var scopes []string
	for key := range flow.Scopes {
		scopes = append(scopes, key)
	}

	var oauthTokenSource oauth2.TokenSource
	var err error
	switch flow.GrantType {
	case discovery.ClientCredentials:
		oauthTokenSource, err = createClientCredentialsTokenSource(ctx, clientId, secret, scopes, flow.TokenURL), nil
		break
	case discovery.Implicit:
		l.V(logger.ImportantLevel).Info("target is using OAuth implicit flow, this flow is deprecated, consider using the authorization code flow instead, https://oauth.net/2/grant-types/implicit/")
		return nil, errors.New("implicit OAuth2 flow is unsupported")
	case discovery.Password:
		oauthTokenSource, err = createPasswordTokenSource(ctx, clientId, secret, scopes, flow.AuthorizationURL, flow.TokenURL)
		break
	case discovery.AuthorizationCode:
		oauthTokenSource, err = createAuthorizationCodeTokenSource(ctx, clientId, secret, scopes, flow.AuthorizationURL, flow.TokenURL)
		break
	}

	if err != nil {
		return nil, fmt.Errorf("error while creating oauth token source: %w", err)
	}
	// Convert to regular token source
	return CreateOAuthTokenSource(oauthTokenSource), nil
}

// createClientCredentialsTokenSource creates an oauth2 ITokenSource for the Client Credentials OAuth flow
// https://tools.ietf.org/html/rfc6749#section-4.4
func createClientCredentialsTokenSource(ctx context.Context, clientId string, clientSecret string, scopes []string, tokenUrl string) oauth2.TokenSource {
	conf := clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     tokenUrl,
		Scopes:       scopes,
		AuthStyle:    oauth2.AuthStyleAutoDetect,
	}
	return conf.TokenSource(ctx)
}

// createPasswordTokenSource creates an oauth2 ITokenSource for the Password OAuth flow
// https://tools.ietf.org/html/rfc6749#section-4.3
func createPasswordTokenSource(ctx context.Context, clientId string, clientSecret string, scopes []string, authUrl string, tokenUrl string) (oauth2.TokenSource, error) {
	conf := oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL:  tokenUrl,
			AuthURL:   authUrl,
			AuthStyle: oauth2.AuthStyleAutoDetect,
		},
		Scopes: scopes,
	}
	tok, err := conf.PasswordCredentialsToken(ctx, clientId, clientSecret)
	if err != nil {
		return nil, err
	}
	return conf.TokenSource(ctx, tok), nil
}

// createAuthorizationCodeTokenSource creates an oauth2 ITokenSource for the Authorization code OAuth flow
// https://tools.ietf.org/html/rfc6749#section-4.1
func createAuthorizationCodeTokenSource(ctx context.Context, clientId string, clientSecret string, scopes []string, authUrl string, tokenUrl string) (oauth2.TokenSource, error) {
	conf := oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL:  tokenUrl,
			AuthURL:   authUrl,
			AuthStyle: oauth2.AuthStyleAutoDetect,
		},
		Scopes: scopes,
	}

	// Offline access enables the use of refresh tokens
	tok, err := conf.Exchange(ctx, clientId, oauth2.AccessTypeOffline)
	if err != nil {
		return nil, err
	}
	return conf.TokenSource(ctx, tok), nil
}
