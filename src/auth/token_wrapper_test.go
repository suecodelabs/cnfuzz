package auth

import (
	"github.com/stretchr/testify/assert"
	"github.com/suecodelabs/cnfuzz/src/discovery"
	"testing"
)

func TestCreateEmptyTokenWrapper(t *testing.T) {
	_, err := CreateTokenWrapper(discovery.SecuritySchema{}, "", "")
	assert.Errorf(t, err, "no tokensource available for  auth scheme")
}

func TestCreateTokenSourceBasic(t *testing.T) {
	clientId := "myclient"
	secret := "mysecret"
	wrapper := TokenWrapper{
		Schema: discovery.SecuritySchema{
			Type: discovery.BasicSecSchemaType,
		},
		ClientId: clientId,
		Secret:   secret,
	}
	err := wrapper.CreateTokenSource()
	assert.NoError(t, err)
	assert.NotNil(t, wrapper.TokenSource)
}

func TestCreateTokenSourceApiKey(t *testing.T) {
	clientId := "myclient"
	secret := "mysecret"
	wrapper := TokenWrapper{
		Schema: discovery.SecuritySchema{
			Type: discovery.ApiKeySecSchemaType,
		},
		ClientId: clientId,
		Secret:   secret,
	}
	err := wrapper.CreateTokenSource()
	assert.NoError(t, err)
	assert.NotNil(t, wrapper.TokenSource)
}
