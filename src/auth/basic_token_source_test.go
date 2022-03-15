package auth

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBasicAuthToken(t *testing.T) {
	clientId := "some-client"
	secret := "some-secret"
	toEncode := fmt.Sprintf("%s:%s", clientId, secret)
	token := base64.StdEncoding.EncodeToString([]byte(toEncode))

	testSource := &basicAuthTokenSource{
		t: &Token{
			AccessToken: token,
			TokenType:   "basic",
		},
	}
	tok, err := testSource.Token()
	if assert.NoError(t, err) {
		assert.Equal(t, token, tok.AccessToken)
		assert.Equal(t, "basic", tok.TokenType)
	}
}

func TestInvalidBasicToken(t *testing.T) {
	testSource := &apiKeyTokenSource{
		t: &Token{
			Expiry: time.Unix(0, 0),
		},
	}
	tok, err := testSource.Token()
	if assert.Errorf(t, err, "failed to create a new token because the current token is invalid and there is no token source") {
		assert.Nil(t, tok)
	}
}

func TestApiBasicTokenSource(t *testing.T) {
	clientId := "some-client"
	secret := "some-secret"
	toEncode := fmt.Sprintf("%s:%s", clientId, secret)
	token := base64.StdEncoding.EncodeToString([]byte(toEncode))

	src, err := BasicAuthTokenSource(clientId, secret)
	if assert.NoError(t, err) {
		switch src := src.(type) {
		case *basicAuthTokenSource:
			assert.Equal(t, token, src.t.AccessToken)
			assert.Equal(t, "basic", src.t.TokenType)
			assert.Nil(t, src.new)
		}
	}
}
