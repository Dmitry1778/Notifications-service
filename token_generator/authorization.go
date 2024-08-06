package token_generator

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"notify/internal/domain"
	"time"
)

const expirationTime = time.Hour * 72

func (tg *TokenGenerator) ParseToken(accessToken string) (int, error) {
	token, err := jwt.NewParser().ParseWithClaims(accessToken, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("fail method")
		}
		return tg.jwtSecretKey, nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}
	username, ok := claims["sub"]
	if !ok {
		return 0, errors.New("invalid token")
	}
	userID, err := tg.userStorage.GetID(context.Background(), username)
	return *userID, nil
}

func (tg *TokenGenerator) Generate(username string) (*Token, error) {
	user, err := tg.userStorage.Get(context.Background(), username)
	if err != nil {
		return nil, err
	}
	claims := jwt.MapClaims{
		"sub": user.Username,
		"exp": time.Now().Add(expirationTime).Unix(),
	}
	return &Token{jwtToken: jwt.NewWithClaims(jwt.SigningMethodHS256, claims), key: tg.jwtSecretKey}, err
}

func New(jwtSecretKey []byte, userStorage userStorage) *TokenGenerator {
	return &TokenGenerator{jwtSecretKey: jwtSecretKey, userStorage: userStorage}
}

type TokenGenerator struct {
	userStorage  userStorage
	jwtSecretKey []byte
}

type userStorage interface {
	Get(ctx context.Context, username string) (*domain.Employee, error)
	GetID(ctx context.Context, username interface{}) (*int, error)
}
