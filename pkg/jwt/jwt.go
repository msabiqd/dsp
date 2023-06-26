package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/SawitProRecruitment/UserService/entity"
	"github.com/SawitProRecruitment/UserService/pkg/apierror"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(id int, phone_number string, privateKey string) (string, error) {
	claim := jwt.MapClaims{
		"id":           id,
		"phone_number": phone_number,
		"expiry_time":  time.Now().Add(30 * time.Minute),
	}

	pemString := privateKey

	block, _ := pem.Decode([]byte(pemString))
	parseResult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	key := parseResult.(*rsa.PrivateKey)

	// rsaKey := key.(*rsa.PrivateKey)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	jwt, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return jwt, nil
}

func VerifyToken(tokenString string, publicKey string) (entity.User, error) {
	pemString := publicKey

	block, _ := pem.Decode([]byte(pemString))
	parseResult, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return entity.User{}, err
	}
	key := parseResult.(*rsa.PublicKey)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return entity.User{}, apierror.New(errors.New("Invalid token"), http.StatusForbidden)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		now := time.Now()
		time, _ := time.Parse(time.RFC3339Nano, claims["expiry_time"].(string))
		if now.After(time) {
			return entity.User{}, apierror.New(errors.New("Token expired"), http.StatusForbidden)
		}
		return entity.User{
			Id:          int(claims["id"].(float64)),
			PhoneNumber: claims["phone_number"].(string),
		}, nil
	}

	return entity.User{}, apierror.New(errors.New("Invalid token"), http.StatusForbidden)

}
