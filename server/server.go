package server

import (
	"github.com/labstack/echo"
	"scoutiq_server/controllers"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/middleware"
	"strings"
)

type Config struct {
	SecretKey string
	ServerName string
}

type Serv struct {
	*echo.Echo
}

func ServerHeader(cfg *Config) echo.MiddlewareFunc {
	return func (next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderServer, cfg.ServerName)
			c.Set("SecretKey", cfg.SecretKey)
			return next(c)
		}
	}
}

func MyMiddleware(database *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("database", database)
			return next(c)
		}
	}
}

func NewServer(database *gorm.DB, cfg *Config) *Serv {
	s := Serv{echo.New()}
	s.Pre(middleware.RemoveTrailingSlash())
	s.Use(middleware.Logger())
	s.Use(middleware.RequestID())

	apiGroup := s.Group("/api/v3")
	apiGroup.Use(MyMiddleware(database))
	apiGroup.Use(ServerHeader(cfg))
	apiGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Skipper: func(c echo.Context) bool {
			if strings.HasPrefix(c.Path(), "/api/v3/token") {
				return true
			}
			return false
		},
		SigningKey: []byte(cfg.SecretKey),
	}))
	controllers.Register(apiGroup)

	return &s
}