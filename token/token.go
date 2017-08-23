package token

import (
	"strings"
	"time"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"errors"
)

type TokenManager struct {
	SigningMethod jwt.SigningMethod
	Timeout time.Duration
	PrivateKeyLoader func() interface{}
	PublicKeyLoader func() interface{}
}

type CreateTokenManagerOption struct {
	// Signing algorithm.
	// Possible values are HS256, HS384, HS512, RS256, RS384, RS512, PS256, PS384, PS512, ES256, ES384, ES512.
	SigningAlgorithm string
	Timeout time.Duration
	PrivateKeyLoader func() interface{}
	PublicKeyLoader func() interface{}
}

func CreateTokenManager(o CreateTokenManagerOption) (*TokenManager, error) {

	method := jwt.GetSigningMethod(o.SigningAlgorithm)
	if method == nil {
		return nil, errors.New("invalid signing algorithm")
	}

	if o.PrivateKeyLoader == nil {
		return nil, errors.New("private key loader required")
	}

	if o.PublicKeyLoader == nil {
		return nil, errors.New("public key loader required")
	}

	m := &TokenManager{
		SigningMethod: method,
		Timeout: o.Timeout,
		PrivateKeyLoader: o.PrivateKeyLoader,
		PublicKeyLoader: o.PublicKeyLoader,
	}

	if m.Timeout == 0 {
		m.Timeout = time.Hour * 1
	}

	return m, nil
}

func (m *TokenManager) CreateToken(username string) *jwt.Token {

	return jwt.NewWithClaims(m.SigningMethod, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(m.Timeout).Unix(),
	})
}

func (m *TokenManager) CreateSignedToken(token *jwt.Token) (string, error) {

	tokenString, err := token.SignedString(m.PrivateKeyLoader())
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func (m *TokenManager) ParseTokenFromRequest(r *http.Request) (*jwt.Token, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return nil, errors.New("Auth header empty")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, errors.New("Invalid auth header")
	}

	return m.ParseToken(parts[1])
}

func (m *TokenManager) ParseToken(signedToken string) (*jwt.Token, error) {
	return jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if m.SigningMethod != token.Method {
			return nil, errors.New("Invalid signing algorithm")
		}
		return m.PublicKeyLoader(), nil
	})
}