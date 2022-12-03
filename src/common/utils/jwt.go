package utils

import (
	"encoding/base64"
	"github.com/golang-jwt/jwt/v4"
	"medusa/src/api/config"
	"time"
)

type Jwt struct {
	cfg config.JwtConfig
}

func NewJwt(cfg config.JwtConfig) *Jwt {
	return &Jwt{
		cfg: cfg,
	}
}

func (j *Jwt) CreateJwt(payload interface{}) (string, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(j.cfg.AccessTokenKey)
	if err != nil {
		return "", err
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()
	claims := make(jwt.MapClaims)
	claims["sub"] = payload
	claims["exp"] = now.Add(time.Duration(j.cfg.Ttl)).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", err
	}

	return token, nil
}
