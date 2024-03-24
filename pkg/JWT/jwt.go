package JWT

import (
	"github.com/golang-jwt/jwt"
	"sso/pkg/apperrors"
	"time"
)

type jwtToken struct {
	signingKey string
	ttl        time.Duration
}

func New(signingKey string, ttl time.Duration) (jwtToken, error) {
	if signingKey == "" {
		return jwtToken{}, apperrors.ErrNoSigningKey
	}

	return jwtToken{signingKey: signingKey, ttl: ttl}, nil
}

// New creates new JWT token with claims and subject in payload
func (j jwtToken) New(sub string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   sub,
		ExpiresAt: time.Now().Add(j.ttl).Unix(),
	})

	return token.SignedString([]byte(j.signingKey))
}

// Parse parses and validating JWT token, returns subject
func (j jwtToken) Parse(token string) (string, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.ErrUnexpectedSignMethod
		}
		return []byte(j.signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok && !t.Valid {
		return "", apperrors.ErrNoClaims
	}
	return claims["sub"].(string), nil
}
