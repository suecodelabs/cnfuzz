package auth

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/suecodelabs/cnfuzz/src/discovery"
	"testing"
)

func TestCreateEmptyTokenWrapper(t *testing.T) {
	_, err := CreateTokenSource(discovery.SecuritySchema{}, "", "")
	assert.Errorf(t, err, "no tokensource available for  auth scheme")
}

func TestCreateTokenSourceBasic(t *testing.T) {
	clientId := "myclient"
	secret := "mysecret"
	tSource, err := CreateTokenSource(discovery.SecuritySchema{
		Type: discovery.BasicSecSchemaType,
	}, clientId, secret)
	assert.NoError(t, err)
	if assert.NotNil(t, tSource) {
		tok, tErr := tSource.Token()
		if assert.NoError(t, tErr) {
			assert.NotNil(t, tok)
			expectedTok := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientId, secret)))
			assert.Equal(t, expectedTok, tok.AccessToken)
			// TODO check header
		}
	}
}

func TestCreateTokenSourceApiKey(t *testing.T) {
	clientId := "myclient"
	secret := "mysecret"
	tSource, err := CreateTokenSource(discovery.SecuritySchema{
		Type: discovery.ApiKeySecSchemaType,
	}, clientId, secret)
	assert.NoError(t, err)
	if assert.NotNil(t, tSource) {
		tok, tErr := tSource.Token()
		if assert.NoError(t, tErr) {
			assert.NotNil(t, tok)
			assert.Equal(t, secret, tok.AccessToken)
			// TODO check header
		}
	}
}

/*
TODO create a test for oauth token source
func TestCreateTokenSourceOAuth(t *testing.T)
*/
