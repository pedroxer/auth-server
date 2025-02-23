package routes

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

type userImpl struct {
	r *Router
}

func registerUserRoutes(r *Router) *userImpl {
	impl := &userImpl{
		r: r,
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
	}
	type loginResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

}

func (impl *userImpl) validateToken(ctx *fasthttp.RequestCtx) {

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

}

func (impl *userImpl) invalidateToken(ctx *fasthttp.RequestCtx) {

}
