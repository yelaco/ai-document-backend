package auth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgon2PasswordHasher_HashAndVerify(t *testing.T) {
	hasher := NewArgon2PasswordHasher()
	password := "thisisasecurePassword123456789@"

	hash, err := hasher.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	err = hasher.VerifyPassword(hash, password)
	assert.NoError(t, err)
}

func TestArgon2PasswordHasher_VerifyIncorrectPassword(t *testing.T) {
	hasher := NewArgon2PasswordHasher()
	password := "thisisasecurePassword123456789@"
	wrongPassword := "wrongPassword"

	hash, err := hasher.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	err = hasher.VerifyPassword(hash, wrongPassword)
	assert.Error(t, err)
}

func TestArgon2PasswordHasher_VerifyMalformedHash(t *testing.T) {
	hasher := NewArgon2PasswordHasher()
	// Malformed hash: too few $ segments
	malformedHash := "$argon2id$v=19$m=65536,t=3,p=2$abcd"
	password := "irrelevant"

	err := hasher.VerifyPassword(malformedHash, password)
	assert.Error(t, err)
}

func TestArgon2PasswordHasher_VerifyParamsMismatch(t *testing.T) {
	hasher := NewArgon2PasswordHasher()
	password := "thisisasecurePassword123456789@"

	hash, err := hasher.HashPassword(password)
	assert.NoError(t, err)

	// Tamper with hash params so they mismatch
	tamperedHash := hash
	tamperedHash = strings.Replace(tamperedHash, "m=65536", "m=32768", 1)

	err = hasher.VerifyPassword(tamperedHash, password)
	assert.Error(t, err)
}
