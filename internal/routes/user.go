package routes

import (
	"encoding/json"

	commonerrors "github.com/pedroxer/auth-service/internal/routes/common_errors"
	"github.com/pedroxer/auth-service/internal/service"
	"github.com/pedroxer/auth-service/internal/service/auth"
	"github.com/valyala/fasthttp"
)


type userImpl struct {
	r    *Router
	auth auth.Auth
}

func registerUserRoutes(r *Router, auth auth.Auth) {
	impl := userImpl{
		r:    r,
		auth: auth,
	}

	r.rtr.POST("/api/v1/login", impl.login)
	r.rtr.GET("/api/v1/validate", impl.validateToken)
	r.rtr.POST("/api/v1/refresh", impl.refreshToken)
	r.rtr.DELETE("/api/v1/invalidate_token", impl.invalidateToken)
}

func (impl *userImpl) login(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	type loginRequest struct {
		Email    string `json:"email"`
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
	access, refresh, err := impl.auth.Login(req.Email, req.Password, req.AppId)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			commonerrors.FormError(ctx, err.Error(), fasthttp.StatusBadRequest)
		case service.ErrRestricted:
			commonerrors.FormError(ctx, err.Error(), fasthttp.StatusForbidden)
		default:
			commonerrors.FormError(ctx, err.Error(), fasthttp.StatusInternalServerError)
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
	ctx.SetContentType("application/json")
	type validateRequest struct {
		AccessToken string `json:"access_token"`
		AppID       int    `json:"app_id"`
	}
	var req validateRequest
	if err := json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	success, err := impl.auth.Validate(req.AccessToken, req.AppID)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			commonerrors.FormError(ctx, err.Error(), fasthttp.StatusBadRequest)
		case service.ErrRestricted:
			commonerrors.FormError(ctx, err.Error(), fasthttp.StatusForbidden)
		case service.ErrInvalidToken:
			commonerrors.FormError(ctx, err.Error(), fasthttp.StatusBadRequest)
		case service.ErrTokenExpired:
			commonerrors.FormError(ctx, err.Error(), fasthttp.StatusUnauthorized)
		default:
			commonerrors.FormError(ctx, err.Error(), fasthttp.StatusInternalServerError)
		}
		return
	}
	body, err := json.Marshal(success)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	
	ctx.SetBody(body)
}

func (impl *userImpl) refreshToken(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	type refreshTokenRequest struct {
		RefreshToken string `json:"refresh_token"`
		AppID        int    `json:"app_id"`
	}

	req := refreshTokenRequest{}
	if err := json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		ctx.SetBodyString(err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	access, err := impl.auth.Refresh(req.RefreshToken, req.AppID)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	body, err := json.Marshal(access)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			commonerrors.FormError(ctx, err.Error(), fasthttp.StatusBadRequest)
		case service.ErrRestricted:
			commonerrors.FormError(ctx, err.Error(), fasthttp.StatusForbidden)
		case service.ErrInvalidToken:
			commonerrors.FormError(ctx, err.Error(), fasthttp.StatusBadRequest)
		case service.ErrTokenExpired:
			commonerrors.FormError(ctx, err.Error(), fasthttp.StatusUnauthorized)
		default:
			commonerrors.FormError(ctx, err.Error(), fasthttp.StatusInternalServerError)
		}
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.SetBody(body)
}

func (impl *userImpl) invalidateToken(ctx *fasthttp.RequestCtx) {

}
