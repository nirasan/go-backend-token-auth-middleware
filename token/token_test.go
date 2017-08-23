package token

import (
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
)

func TestCreateTokenManager(t *testing.T) {
	m, err := CreateTokenManager(CreateTokenManagerOption{})
	if err == nil {
		t.Error("SignedMethod required")
	}

	m, err = CreateTokenManager(CreateTokenManagerOption{
		SigningAlgorithm: "ES256",
		PrivateKeyLoader: LoadPrivateKey,
		PublicKeyLoader:  LoadPublicKey,
	})
	if err != nil {
		t.Error(err)
	}
	if m.Timeout != time.Hour {
		t.Errorf("invalid default timeout: %+v", m)
	}
}

func TestTokenManager_CreateToken(t *testing.T) {
	m, _ := CreateTokenManager(CreateTokenManagerOption{
		SigningAlgorithm: "ES256",
		PrivateKeyLoader: LoadPrivateKey,
		PublicKeyLoader:  LoadPublicKey,
	})

	token := m.CreateToken("user1")

	if claims, ok := token.Claims.(jwt.MapClaims); !ok || claims["sub"] != "user1" {
		t.Errorf("invalid token %+v", token)
	}
}

func TestTokenManager_CreateSignedToken(t *testing.T) {
	m, _ := CreateTokenManager(CreateTokenManagerOption{
		SigningAlgorithm: "ES256",
		PrivateKeyLoader: LoadPrivateKey,
		PublicKeyLoader:  LoadPublicKey,
	})

	signedToken, err := m.CreateSignedToken(m.CreateToken("user1"))
	if err != nil {
		t.Error(err)
	}

	t.Logf("%q", signedToken)

	token, err := m.ParseToken(signedToken)

	t.Logf("%q", token)

	if claims, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid || claims["sub"] != "user1" {
		t.Errorf("invalid token %+v", token)
	}
}

func LoadPrivateKey() interface{} {
	keyData, e := ioutil.ReadFile("./test/ec256-private.pem")
	if e != nil {
		panic(e.Error())
	}
	key, e := jwt.ParseECPrivateKeyFromPEM(keyData)
	if e != nil {
		panic(e.Error())
	}
	return key
}

func LoadPublicKey() interface{} {
	keyData, e := ioutil.ReadFile("./test/ec256-public.pem")
	if e != nil {
		panic(e.Error())
	}
	key, e := jwt.ParseECPublicKeyFromPEM(keyData)
	if e != nil {
		panic(e.Error())
	}
	return key
}
