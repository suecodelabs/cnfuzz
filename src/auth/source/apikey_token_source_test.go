package source

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestApiKeyToken(t *testing.T) {
	secret := "some-secret"
	testSource := &apiKeyTokenSource{
		t: &Token{
			AccessToken: secret,
			TokenType:   "api-key",
		},
	}
	tok, err := testSource.Token()
	if assert.NoError(t, err) {
		assert.Equal(t, secret, tok.AccessToken)
		assert.Equal(t, "api-key", tok.TokenType)
	}
}

func TestInvalidApiKeyToken(t *testing.T) {
	secret := "some-secret"
	testSource := &apiKeyTokenSource{
		t: &Token{
			Expiry:      time.Unix(0, 0),
			AccessToken: secret,
		},
	}
	tok, err := testSource.Token()
	if assert.Errorf(t, err, "failed to create a new api key token because the current token is invalid and there is no token source") {
		assert.Nil(t, tok)
	}
}

func TestApiKeyTokenSource(t *testing.T) {
	secret := "some-secret"
	src, err := ApiKeyTokenSource(secret)
	if assert.NoError(t, err) {
		switch src := src.(type) {
		case *apiKeyTokenSource:
			assert.Equal(t, secret, src.t.AccessToken)
			assert.Equal(t, "api-key", src.t.TokenType)
			assert.Nil(t, src.new)
		}
	}
}
