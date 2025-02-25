package auth

import (
	"github.com/pedroxer/auth-service/internal/models"
	"github.com/pedroxer/auth-service/internal/pkg/hashing"
	"github.com/pedroxer/auth-service/internal/pkg/jwt"
	"github.com/pedroxer/auth-service/internal/service"
	log "github.com/sirupsen/logrus"
)

type UserProvider interface {
	GetUser(email string) (models.User, error)
	GetUserByID(userID int) (models.User, error)
}
type AppProvider interface {
	AddApp(app models.App) error
	GetApp(appId int) (models.App, error)
	GetAppSecret(appID int) (string, error)
}
type UserAuth struct {
	UsPrv  UserProvider
	AppPrv AppProvider
	logger *log.Logger
}

func NewUserAuth(userProv UserProvider, appProv AppProvider, logger *log.Logger) *UserAuth {
	return &UserAuth{
		UsPrv:  userProv,
		AppPrv: appProv,
		logger: logger,
	}
}

func (ua *UserAuth) Validate(token string, appId int) (bool, error) {
	appSecret, err := ua.AppPrv.GetAppSecret(appId)
	if err != nil {
		return false, err
	}
	payload, err := jwt.ValidateToken(token, appSecret)
	if err != nil {
		return false, err
	}
	if payload.IsRefresh {
		return false, service.ErrInvalidToken
	}
	return true, nil
}

func (ua *UserAuth) Login(username, password string, appId int) (string, string, error) {
	user, err := ua.UsPrv.GetUser(username)
	if err != nil {
		return "", "", err
	}

	app, err := ua.AppPrv.GetApp(appId)
	if err != nil {
		return "", "", err
	}
	err = hashing.CheckPassword(password, user.Password)
	if err == nil {
		return "", "", service.ErrUserNotFound
	}
	access, refresh, err := jwt.GenerateTokens(app, user)
	if err != nil {
		return "", "", err
	}
	ua.logger.Infof("user %d logged in", user.ID)
	return access, refresh, nil
}

func (ua *UserAuth) Refresh(refreshToken string, appID int) (string, error) {
	app, err := ua.AppPrv.GetApp(appID)
	if err != nil {
		return "", err
	}
	payload, err := jwt.ValidateToken(refreshToken, app.Secret)
	if err != nil {
		return "", err
	}
	if !payload.IsRefresh {
		return "", service.ErrInvalidToken
	}
	user, err := ua.UsPrv.GetUserByID(payload.UserID)
	if err != nil {
		return "", err
	}
	newAccessToken, err := jwt.NewAccessToken(app, user)
	if err != nil {
		return "", err
	}
	return newAccessToken, nil
}
