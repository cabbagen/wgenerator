package providers

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const JWTTokenSecret string = "gateway-token-secret"

type JWTTokenClaims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

func HandleGenerateJWTToken(userId string) (string, error) {
	claims := JWTTokenClaims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(JWTTokenSecret))
}

func HandleParseJWTToken(tokenString string) (string, error) {
	token, error := jwt.ParseWithClaims(tokenString, &JWTTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTTokenSecret), nil
	})

	if error != nil {
		return "", error
	}
	if claims, ok := token.Claims.(*JWTTokenClaims); ok {
		return claims.UserId, nil
	}
	return "", fmt.Errorf("unknown claims type, cannot proceed")
}
