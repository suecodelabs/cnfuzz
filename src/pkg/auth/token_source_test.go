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
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/suecodelabs/cnfuzz/src/pkg/discovery"
	"github.com/suecodelabs/cnfuzz/src/pkg/logger"
	"testing"
)

func TestCreateEmptyTokenWrapper(t *testing.T) {
	l := logger.CreateDebugLogger()
	_, err := CreateTokenSource(l, discovery.SecuritySchema{}, "", "")
	assert.Errorf(t, err, "no tokensource available for  auth scheme")
}

func TestCreateTokenSourceBasic(t *testing.T) {
	l := logger.CreateDebugLogger()
	clientId := "myclient"
	secret := "mysecret"
	tSource, err := CreateTokenSource(l, discovery.SecuritySchema{
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
	l := logger.CreateDebugLogger()
	clientId := "myclient"
	secret := "mysecret"
	tSource, err := CreateTokenSource(l, discovery.SecuritySchema{
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
