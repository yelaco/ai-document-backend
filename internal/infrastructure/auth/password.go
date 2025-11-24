package auth

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/yelaco/ai-document-backend/pkg/util"
	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = fmt.Errorf("invalid hash format")
	ErrIncompatibleVersion = fmt.Errorf("incompatible argon2 version")
	ErrParamsMismatch      = fmt.Errorf("argon2 parameters mismatch")
)

type Argon2PasswordHasher struct {
	params params
}

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func NewArgon2PasswordHasher() util.PasswordHasher {
	return &Argon2PasswordHasher{
		params: params{
			memory:      64 * 1024,
			iterations:  3,
			parallelism: 2,
			saltLength:  16,
			keyLength:   32,
		},
	}
}

func (ph *Argon2PasswordHasher) HashPassword(password string) (string, error) {
	salt, err := util.GenerateRandomBytes(int(ph.params.saltLength))
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	hash := argon2.IDKey([]byte(password), salt, ph.params.iterations, ph.params.memory, ph.params.parallelism, ph.params.keyLength)

	// Base64 encode the salt and hash for storage
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Standared encoded hash representation
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, ph.params.memory, ph.params.iterations, ph.params.parallelism, b64Salt, b64Hash)
	return encodedHash, nil
}

func (ph *Argon2PasswordHasher) VerifyPassword(hashedPassword, password string) error {
	salt, hash, err := ph.decodeHash(hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to decode hash: %w", err)
	}

	otherHash := argon2.IDKey([]byte(password), salt, ph.params.iterations, ph.params.memory, ph.params.parallelism, ph.params.keyLength)
	if err != nil {
		return fmt.Errorf("failed to hash password for verification: %w", err)
	}

	if subtle.ConstantTimeCompare(hash, otherHash) != 1 {
		return fmt.Errorf("password does not match")
	}
	return nil
}

func (ph *Argon2PasswordHasher) decodeHash(encodedHash string) (salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, ErrIncompatibleVersion
	}

	p := params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	if !compareParams(ph.params, p) {
		return nil, nil, ErrParamsMismatch
	}

	return salt, hash, nil
}

func compareParams(p1, p2 params) bool {
	return p1.memory == p2.memory &&
		p1.iterations == p2.iterations &&
		p1.parallelism == p2.parallelism &&
		p1.saltLength == p2.saltLength &&
		p1.keyLength == p2.keyLength
}
