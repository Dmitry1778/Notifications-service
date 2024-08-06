package token_generator

import "github.com/golang-jwt/jwt/v5"

func (t *Token) String() (string, error) {
	return t.jwtToken.SignedString(t.key)
}

type Token struct {
	jwtToken *jwt.Token
	key      []byte
}
