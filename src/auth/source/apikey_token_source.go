package source

import (
	"errors"
)

type apiKeyTokenSource struct {
	new TokenSource // called when t is expired.
	t   *Token
}

func (s *apiKeyTokenSource) Token() (*Token, error) {
	if s.t.Valid() {
		return s.t, nil
	}
	return nil, errors.New("failed to create a new api key token because the current token is invalid and there is no token source")
}

func ApiKeyTokenSource(secret string) (TokenSource, error) {
	if len(secret) == 0 {
		return nil, errors.New("failed to create token because secret is empty")
	}
	token := &Token{
		AccessToken: secret,
		TokenType:   "api-key",
	}

	return &apiKeyTokenSource{
		t: token,
		// new: tkr,
	}, nil
}
