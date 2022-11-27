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
