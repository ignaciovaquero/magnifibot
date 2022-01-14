package utils

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// ValidateTokenFromHeader validates a JWT token from an Authorization header, with
// the secret used to sign the token.
func ValidateTokenFromHeader(header, secret string) error {
	if secret == "" {
		return nil
	}

	splitToken := strings.Split(header, "Bearer")

	if len(splitToken) != 2 {
		return fmt.Errorf("error getting token from Authorization header: header is not in proper format")
	}

	tok, err := jwt.Parse(strings.Trim(splitToken[1], " "), func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return fmt.Errorf("error parsing token: %w", err)
	}

	if !tok.Valid {
		return fmt.Errorf("token is invalid")
	}

	return nil
}
