// Copyright 2022 Sue B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/suecodelabs/cnfuzz/src/log"
)

const expiryDelta = 10 * time.Second

// Token object holding token information
type Token struct {
	AccessToken  string
	TokenType    string
	RefreshToken string
	Expiry       time.Time
}

// Valid checks if this token is still valid
// checks if this token holds a value and isn't expired yet
func (t *Token) Valid() bool {
	return t != nil && t.AccessToken != "" && !t.expired()
}

// expired checks if this token is expired
func (t *Token) expired() bool {
	if t.Expiry.IsZero() {
		return false
	}
	return t.Expiry.Round(0).Add(-expiryDelta).Before(time.Now())
}

// Type returns the type of this token
// formats the TokenType to a value that can be used in the Authorization http header
func (t *Token) Type() string {
	log.L().Debugf("using token type: %s (this is also used in token header Authorizations: <token type prefix> <tokem>", t.TokenType)
	if strings.EqualFold(t.TokenType, "bearer") {
		return "Bearer"
	}
	if strings.EqualFold(t.TokenType, "mac") {
		return "MAC"
	}
	if strings.EqualFold(t.TokenType, "basic") {
		return "Basic"
	}
	if strings.EqualFold(t.TokenType, "api-key") {
		return ""
	}
	if t.TokenType != "" {
		return t.TokenType
	}
	return "Bearer"
}

// SetAuthHeader set the authorization header in a http Request using this token
func (t *Token) SetAuthHeader(r *http.Request) {
	r.Header.Set("Authorization", t.CreateAuthHeaderValue())
}

// CreateAuthHeader returns a http Header holding the Authorization header from this token
func (t *Token) CreateAuthHeader() http.Header {
	return http.Header{
		"Authorization": {t.CreateAuthHeaderValue()},
	}
}

// CreateAuthHeaderValue creates the value for the Authorization header from this token
func (t *Token) CreateAuthHeaderValue() string {
	return t.Type() + " " + t.AccessToken
}
