package auth

import (
	"encoding/base64"
	"errors"
)

// basicAuthTokenSource ITokenSource for basic authentication
// https://datatracker.ietf.org/doc/html/rfc7617
type basicAuthTokenSource struct {
	new ITokenSource // called when t is expired.
	t   *Token
}

// Token create a new basic auth token
func (s *basicAuthTokenSource) Token() (*Token, error) {
	if s.t.Valid() {
		return s.t, nil
	}
	return nil, errors.New("failed to create a new basic token because the current token is invalid and there is no token source")
}

// BasicAuthTokenSource creates a new basicAuthTokenSource for basic authentication
func BasicAuthTokenSource(clientId string, secret string) (ITokenSource, error) {

	if len(secret) == 0 {
		return nil, errors.New("failed to create token because secret is empty")
	}
	toEncode := clientId + ":" + secret
	aToken := base64.StdEncoding.EncodeToString([]byte(toEncode))

	token := &Token{
		AccessToken: aToken,
		TokenType:   "basic",
	}

	return &basicAuthTokenSource{
		t: token,
		// new: tkr,
	}, nil
}
