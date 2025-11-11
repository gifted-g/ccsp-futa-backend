package test

import (
	"testing"
	
	// Assuming the code being tested is in the 'auth' package of your project
	auth "ccsp-futa-alumni/handlers" 

	"github.com/stretchr/testify/assert"
)

func TestGenerateOTP(t *testing.T) {
	// Must use the exported name (GenerateOTP) and package prefix (auth.)
	otp := auth.GenerateOTP() 
	assert.Equal(t, 6, len(otp))
}

func TestHashPasswordAndCompare(t *testing.T) {
	pwd := "mysecret"
	// Must use the exported name (HashPassword) and package prefix (auth.)
	hashed, err := auth.HashPassword(pwd) 
	assert.NoError(t, err)
	// Must use the exported name (CheckPasswordHash) and package prefix (auth.)
	assert.True(t, auth.CheckPasswordHash(pwd, hashed))
}