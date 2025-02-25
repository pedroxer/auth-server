package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/pedroxer/auth-service/internal/models"
	payload2 "github.com/pedroxer/auth-service/internal/pkg/payload"
	"github.com/pedroxer/auth-service/internal/service"
	"time"
)

func NewAccessToken(app models.App, user models.User) (string, error) {
	payload := payload2.NewAccessToken(app, user)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, payload)
	accessTokenString, err := accessToken.SignedString([]byte(app.Secret))

	return accessTokenString, err
}

func NewRefreshToken(app models.App, user models.User) (string, error) {
	payload := payload2.NewRefreshToken(app, user)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, payload)
	refreshTokenString, err := refreshToken.SignedString([]byte(app.Secret))

	return refreshTokenString, err
}

func GenerateTokens(app models.App, user models.User) (string, string, error) {
	accessToken, err := NewAccessToken(app, user)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := NewRefreshToken(app, user)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func ValidateToken(tokenString string, appSecret string) (*payload2.Payload, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&payload2.Payload{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(appSecret), nil
		},
		jwt.WithLeeway(2*time.Minute),
		jwt.WithValidMethods([]string{"HS512"}),
	)

	if err != nil {
		return nil, err // Обработка ошибок как в предыдущем ответе
	}

	claims, ok := token.Claims.(*payload2.Payload)
	if !ok || !token.Valid {
		return nil, service.ErrInvalidToken
	}

	// Проверка срока действия через RegisteredClaims
	if claims.ExpiresAt.Before(time.Now()) {
		return nil, service.ErrTokenExpired
	}
	return claims, nil
}
