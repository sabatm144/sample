package authentication

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

//JWTAuthenticationBackend provides properties of authentication entity
type JWTAuthenticationBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

const (
	tokenDuration = 72
	expireOffset  = 3600
)

type AuthStruct struct {
	ID         string
	IsCustomer bool
}

var authBackendInstance *JWTAuthenticationBackend

//InitJWTAuthenticationBackend inits the authentication key
func InitJWTAuthenticationBackend() *JWTAuthenticationBackend {
	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthenticationBackend{
			privateKey: getPrivateKey(),
			PublicKey:  getPublicKey(),
		}
	}

	return authBackendInstance
}

//GetAuthMethod returns the func for extracting a token
func GetAuthMethod() jwt.Keyfunc {
	backend := InitJWTAuthenticationBackend()
	signingKey, _ := json.Marshal(backend.PublicKey)

	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("invalid token")
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else {
			// log.Println("valid token")
			return signingKey, nil
		}

	}
}

//GenerateToken generates token for a user
func (backend *JWTAuthenticationBackend) GenerateToken(user AuthStruct) (string, error) {

	var signingKey []byte
	signingKey, err := json.Marshal(backend.PublicKey)

	claims := jwt.MapClaims{
		"exp":        time.Now().Add(time.Hour * time.Duration(tokenDuration)).Unix(),
		"iat":        time.Now().Unix(),
		"id":         user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signingKey)

	return ss, err
}

func getPrivateKey() *rsa.PrivateKey {
	privateKeyFile, err := os.Open("server/conf/app.rsa")
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	privateKeyFile.Close()

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	return privateKeyImported
}

func getPublicKey() *rsa.PublicKey {
	publicKeyFile, err := os.Open("server/conf/app.rsa.pub")
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := publicKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	publicKeyFile.Close()

	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

	if !ok {
		panic(err)
	}

	return rsaPub
}
