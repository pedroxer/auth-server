package routes

import (
	"time"

	fasthttprouter "github.com/fasthttp/router"
	"github.com/pedroxer/auth-service/internal/config"
	"github.com/pedroxer/auth-service/internal/service/auth"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type Router struct {
	rtr     *fasthttprouter.Router
	srv     *fasthttp.Server
	logger  *logrus.Logger
	port    string
}

func New(cfg *config.Config, logger *logrus.Logger, auth auth.Auth) *Router {

	rtr := fasthttprouter.New()

	r := &Router{
		rtr:     rtr,
		srv: &fasthttp.Server{
			Handler:            rtr.Handler,
			MaxRequestBodySize: 100_000_000,
			ReadTimeout:        time.Duration(cfg.Api.ReadTimeout) * time.Second,
			WriteTimeout:       time.Duration(cfg.Api.WriteTimeout) * time.Second,
			IdleTimeout:        time.Duration(cfg.Api.IdleTimeout) * time.Second,
			Logger:             logger,
		},
		logger: logger,
		port:   cfg.Api.Port,
	}
	registerUserRoutes(r,auth)
	

	r.rtr.HandleMethodNotAllowed = true
	
	return r
}

func (r *Router) Start() error {
	return r.srv.ListenAndServe(r.port)
}
func (r *Router) Shutdown() error {
	return r.srv.Shutdown()
}
