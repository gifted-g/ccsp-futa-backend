package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateOTP(t *testing.T) {
	otp := generateOTP()
	assert.Equal(t, 6, len(otp))
}

func TestHashPasswordAndCompare(t *testing.T) {
	pwd := "mysecret"
	hashed, err := hashPassword(pwd)
	assert.NoError(t, err)
	assert.True(t, checkPasswordHash(pwd, hashed))
}
