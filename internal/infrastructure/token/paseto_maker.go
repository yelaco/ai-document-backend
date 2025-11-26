package token

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

// PasetoMaker is a PASETO token maker
type PasetoV4Maker struct {
	secretKey paseto.V4AsymmetricSecretKey
	publicKey paseto.V4AsymmetricPublicKey
}

// NewPasetoV4Maker creates a new PasetoMaker
func NewPasetoV4Maker() (Maker, error) {
	secretKey := paseto.NewV4AsymmetricSecretKey()
	publicKey := secretKey.Public()

	return &PasetoV4Maker{
		secretKey: secretKey,
		publicKey: publicKey,
	}, nil
}

func (maker *PasetoV4Maker) GetPublicKey() string {
	return maker.publicKey.ExportHex()
}

// CreateToken creates a new token for a specific username and duration
func (maker *PasetoV4Maker) CreateToken(sub, email, role string) (string, error) {
	now := time.Now()
	exp := now.Add(AccessTokenExpirationMinutes * time.Minute)
	nbf := now

	token := paseto.NewToken()
	token.SetAudience(TokenAudience)
	token.SetJti(uuid.NewString())
	token.SetIssuer(TokenIssuer)
	token.SetSubject(sub)

	token.SetExpiration(exp)
	token.SetNotBefore(nbf)
	token.SetIssuedAt(now)

	_ = token.Set("email", email)
	_ = token.Set("role", role)

	signed := token.V4Sign(maker.secretKey, nil)
	return signed, nil
}

// VerifyToken checks if the token is valid or not
func (maker *PasetoV4Maker) VerifyToken(signed string, skipExp bool) (*Claims, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.ForAudience(TokenAudience))
	parser.AddRule(paseto.IssuedBy(TokenIssuer))
	if !skipExp {
		parser.AddRule(paseto.NotExpired())
		parser.AddRule(paseto.ValidAt(time.Now()))
	}

	parsedToken, err := parser.ParseV4Public(maker.publicKey, signed, nil)
	if err != nil {
		return nil, fmt.Errorf("token.PasetoV4Maker.VerifyToken: %w", err)
	}

	sub, err := parsedToken.GetSubject()
	if err != nil {
		return nil, fmt.Errorf("token.PasetoV4Maker.VerifyToken: %w: sub", ErrInvalidToken)
	}

	iss, err := parsedToken.GetIssuer()
	if err != nil {
		return nil, fmt.Errorf("token.PasetoV4Maker.VerifyToken: %w: iss", ErrInvalidToken)
	}

	exp, err := parsedToken.GetExpiration()
	if err != nil {
		return nil, fmt.Errorf("token.PasetoV4Maker.VerifyToken: %w: exp", ErrInvalidToken)
	}

	var email string
	if err = parsedToken.Get("email", &email); err != nil {
		return nil, fmt.Errorf("token.PasetoV4Maker.VerifyToken: %w: email", ErrInvalidToken)
	}

	var role string
	if err = parsedToken.Get("role", &role); err != nil {
		return nil, fmt.Errorf("token.PasetoV4Maker.VerifyToken: %w: role", ErrInvalidToken)
	}

	return &Claims{
		Sub:   sub,
		Email: email,
		Role:  role,
		Iss:   iss,
		Exp:   exp,
	}, nil
}
