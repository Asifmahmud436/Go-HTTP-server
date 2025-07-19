package auth

import (
	"errors"
	"strings"
	"net/http"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pass), err
}

func CheckPassword(password, hash string) (string, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return "", errors.New("wrong password")
	}
	return password, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Issuer: "chirpy", IssuedAt: jwt.NewNumericDate(time.Now()), ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)), Subject: userID.String()})
	result, err := token.SignedString(tokenSecret)
	if err != nil {
		return "", err
	}
	return result, err
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error){
	type MyCustomClaims struct{
		Foo string `json:"foo"`
		jwt.RegisteredClaims
	}
	claims := &MyCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString,claims,func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret),nil
	})
	if err!=nil{
		return uuid.Nil,err
	}
	if !token.Valid{
		return uuid.Nil, errors.New("invalid token")
	}
	givenUser := claims.Subject
	result,err := uuid.Parse(givenUser)
	if err!=nil{
		return uuid.Nil, err
	}
	return result,nil
}

func GetBearerToken(headers http.Header) (string, error){
	auth := headers.Get("Authorization")
	if auth==""{
		return "",errors.New("empty header")
	}
	token := strings.TrimSpace(strings.TrimPrefix(auth,"Bearer "))
	return token,nil
}