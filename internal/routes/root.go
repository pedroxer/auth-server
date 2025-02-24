package routes

import (
	fasthttprouter "github.com/fasthttp/router"
	"github.com/pedroxer/auth-service/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"time"
)

type Router struct {
	rtr    *fasthttprouter.Router
	srv    *fasthttp.Server
	logger *logrus.Logger
	port   string
}

func NewRouter(logger *logrus.Logger, cfg *config.Api, auth Auth) *Router {
	rtr := fasthttprouter.New()
	r := &Router{
		rtr: rtr,
		srv: &fasthttp.Server{
			MaxRequestBodySize: 100_000_000,
			ReadTimeout:        time.Duration(cfg.ReadTimeout) * time.Second,
			WriteTimeout:       time.Duration(cfg.WriteTimeout) * time.Second,
			IdleTimeout:        time.Duration(cfg.IdleTimeout) * time.Second,
			Logger:             logger,
		},
		logger: logger,
		port:   cfg.Port,
	}
	registerUserRoutes(r, auth)
	return r
}
