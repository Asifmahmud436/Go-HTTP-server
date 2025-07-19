package auth_test

import (
	"testing"
	"time"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/Asifmahmud436/Go-HTTP-server/internal/auth"
)

func TestJWT_CreateAndValidate(t *testing.T) {
	secret := "mysecretkey"
	userID := uuid.New()
	expiresIn := time.Minute * 5

	token, err := auth.MakeJWT(userID, secret, expiresIn)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedID, err := auth.ValidateJWT(token, secret)
	assert.NoError(t, err)
	assert.Equal(t, userID, parsedID)
}

func TestJWT_ExpiredToken(t *testing.T) {
	secret := "mysecretkey"
	userID := uuid.New()

	// Token already expired
	token, err := auth.MakeJWT(userID, secret, -1*time.Minute)
	assert.NoError(t, err)

	_, err = auth.ValidateJWT(token, secret)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestJWT_WrongSecret(t *testing.T) {
	secret := "mysecretkey"
	wrongSecret := "wrongsecret"
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, secret, time.Minute*5)
	assert.NoError(t, err)

	_, err = auth.ValidateJWT(token, wrongSecret)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "signature is invalid")
}
