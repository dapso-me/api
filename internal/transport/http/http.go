package http

import (
	"context"
	"net/http"

	_ "api/docs"
	"api/internal/transport/http/helper"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title        Dapso API Documentation
// @version      1.0
// @description  no description
// @BasePath     /api/v1
// @Accept       json
// @Produce      json

type httpServer struct {
	e *echo.Echo
}

func New() *httpServer {
	return &httpServer{
		e: echo.New(),
	}
}

func (h *httpServer) Run(addr string) error {
	return h.e.Start(addr)
}

func (h *httpServer) Shutdown(c context.Context) error {
	return h.e.Shutdown(c)
}

func (h *httpServer) SetupRouter(handlers ...helper.Handler) {
	h.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	h.e.GET("/swagger/*", echoSwagger.WrapHandler)

	baseUrl := h.e.Group("/api/v1")

	for _, handler := range handlers {
		handler.Register(baseUrl, h)
	}
}
