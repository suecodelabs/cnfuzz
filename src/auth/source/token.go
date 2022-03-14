package source

import (
	"net/http"
	"strings"
	"time"

	"github.com/suecodelabs/cnfuzz/src/log"
)

const expiryDelta = 10 * time.Second

type TokenSource interface {
	Token() (*Token, error)
}

type Token struct {
	AccessToken  string
	TokenType    string
	RefreshToken string
	Expiry       time.Time
}

func (t *Token) Valid() bool {
	return t != nil && t.AccessToken != "" && !t.expired()
}

func (t *Token) expired() bool {
	if t.Expiry.IsZero() {
		return false
	}
	return t.Expiry.Round(0).Add(-expiryDelta).Before(time.Now())
}

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

func (t *Token) SetAuthHeader(r *http.Request) {
	r.Header.Set("Authorization", t.CreateAuthHeaderValue())
}

func (t *Token) CreateAuthHeader() http.Header {
	return http.Header{
		"Authorization": {t.CreateAuthHeaderValue()},
	}
}

func (t *Token) CreateAuthHeaderValue() string {
	return t.Type() + " " + t.AccessToken
}
