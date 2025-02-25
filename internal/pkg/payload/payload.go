package payload

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pedroxer/auth-service/internal/models"
	"time"
)

type Payload struct {
	AppID     int    `json:"app_id"`
	Email     string `json:"email"`
	UserID    int    `json:"user_id"`
	Role      int    `json:"role"`
	Position  int    `json:"position"`
	Team      string `json:"team"`
	IsRefresh bool   `json:"is_refresh"` // Флаг для refresh token
	JTI       string `json:"jti"`        // Уникальный идентификатор токена
	jwt.RegisteredClaims
}

// Access Token (TTL 15 минут)
func NewAccessToken(app models.App, user models.User) *Payload {
	return &Payload{
		AppID:     app.ID,
		Email:     user.Email,
		UserID:    user.ID,
		Role:      user.Role,
		Position:  user.Position,
		Team:      user.Team,
		IsRefresh: false,
		JTI:       uuid.New().String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    app.Name,
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}
}

// Refresh Token (TTL 7 дней)
func NewRefreshToken(app models.App, user models.User) *Payload {
	return &Payload{
		AppID:     app.ID,
		Email:     user.Email,
		UserID:    user.ID,
		IsRefresh: true, // Отметка refresh токена
		JTI:       uuid.New().String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			Issuer:    app.Name,
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}
}
