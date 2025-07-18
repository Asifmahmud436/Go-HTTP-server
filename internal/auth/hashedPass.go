package auth

import (
	"golang.org/x/crypto/bcrypt"
	"errors"
)

func HashPassword(password string) (string, error) {
	pass,err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil{
		return "", err
	}
	return string(pass),err
}

func CheckPassword(password, hash string) (string, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil{
		return "",errors.New("wrong password")
	}
	return password, nil
}