package auth

import (
	"fmt"
	"golang.org/x/oauth2"
)

// oAuthTokenSource wraps oauth2.TokenSource
type oAuthTokenSource struct {
	oauthSource oauth2.TokenSource
}

// CreateOAuthTokenSource creates an oAuthTokenSource from a TokenSource from the oauth2 library
func CreateOAuthTokenSource(source oauth2.TokenSource) *oAuthTokenSource {
	return &oAuthTokenSource{
		oauthSource: source,
	}
}

// Token returns a authorization token created using the oauth2 library
func (s *oAuthTokenSource) Token() (*Token, error) {
	oauth2Token, err := s.oauthSource.Token()
	if err != nil {
		return nil, fmt.Errorf("error while getting token from oauth token source: %w", err)
	}
	return tokenFromOAuth2Token(oauth2Token), nil
}

// tokenFromOAuth2Token creates a cnfuzz Token from a oauth2 Token
func tokenFromOAuth2Token(oToken *oauth2.Token) *Token {
	return &Token{
		AccessToken:  oToken.AccessToken,
		TokenType:    oToken.TokenType,
		RefreshToken: oToken.RefreshToken,
		Expiry:       oToken.Expiry,
	}
}
