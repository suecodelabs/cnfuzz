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

package model

import "errors"

var tokens = []Token{}

type Token struct {
	Id    int    `json:"id" example:"1" format:"int32"`
	Token string `json:"token" example:"Token for auth"`
}

func GetToken(token string) (Token, error) {
	for _, storedToken := range tokens {
		if storedToken.Token == token {
			return storedToken, nil
		}
	}
	return Token{}, errors.New("token doesn't exist")
}

func SearchToken(tokenVal string) (token Token, found bool) {
	for _, storedToken := range tokens {
		if storedToken.Token == tokenVal {
			return storedToken, true
		}
	}
	return Token{}, false
}

func CreateToken(token string) int {
	lastId := 0
	if len(tokens) > 0 {
		lastId = tokens[len(tokens)-1].Id
	}
	newToken := Token{
		Id:    lastId + 1,
		Token: token,
	}
	tokens = append(tokens, newToken)
	return newToken.Id
}
