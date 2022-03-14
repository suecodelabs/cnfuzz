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
