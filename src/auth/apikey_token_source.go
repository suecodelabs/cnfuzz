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

import "errors"

// apiKeyTokenSource ITokenSource for API key authentication
type apiKeyTokenSource struct {
	new ITokenSource // called when t is expired.
	t   *Token
}

// Token create a new API key auth token
func (s *apiKeyTokenSource) Token() (*Token, error) {
	if s.t.Valid() {
		return s.t, nil
	}
	return nil, errors.New("failed to create a new api key token because the current token is invalid and there is no token source")
}

// ApiKeyTokenSource creates a new apiKeyTokenSource for API key authentication
func ApiKeyTokenSource(secret string) (ITokenSource, error) {
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
