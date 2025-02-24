package routes

import (
	"encoding/json"
	"github.com/pedroxer/auth-service/internal/service/auth"
	"github.com/valyala/fasthttp"
)

type Auth interface {
	Login(username, password string, appId int) (string, string, error)
	Validate(accessToken string) (bool, error)
	Refresh(refreshToken string) (string, error)
}

type userImpl struct {
	r    *Router
	auth Auth
}

func registerUserRoutes(r *Router, auth Auth) *userImpl {
	impl := &userImpl{
		r:    r,
		auth: auth,
	}

	impl.r.rtr.POST("/api/v1/login", impl.login)
	impl.r.rtr.GET("/api/v1/login", impl.validateToken)
	impl.r.rtr.POST("/api/v1/refresh", impl.refreshToken)
	impl.r.rtr.DELETE("/api/v1/invalidate_token", impl.invalidateToken)
	return impl
}

func (impl *userImpl) login(ctx *fasthttp.RequestCtx) {
	type loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		AppId    int    `json:"app_id"`
	}
	type loginResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	var req loginRequest
	if err := json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	access, refresh, err := impl.auth.Login(req.Username, req.Password, req.AppId)
	if err != nil {
		switch err {
		case auth.ErrUserNotFound:
			ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		case auth.ErrRestricted:
			ctx.Error(err.Error(), fasthttp.StatusForbidden)
		default:
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		}
		return
	}

	resp := loginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}
	body, err := json.Marshal(resp)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(body)
}

func (impl *userImpl) validateToken(ctx *fasthttp.RequestCtx) {
	type validateRequest struct {
		AccessToken string `json:"access_token"`
	}
	var req validateRequest
	if err := json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	success, err := impl.auth.Validate(req.AccessToken)
	if err != nil {
		switch err {
		case auth.ErrTokenExpired:
			ctx.Error(err.Error(), fasthttp.StatusUnauthorized)
		case auth.ErrInvalidToken:
			ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		default:
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		}
		return
	}
	body, err := json.Marshal(success)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(body)
}

func (impl *userImpl) refreshToken(ctx *fasthttp.RequestCtx) {
	type refreshTokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	req := refreshTokenRequest{}
	if err := json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	access, err := impl.auth.Refresh(req.RefreshToken)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	body, err := json.Marshal(access)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(body)
}

func (impl *userImpl) invalidateToken(ctx *fasthttp.RequestCtx) {

}
