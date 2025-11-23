package http

import (
	"context"
	"net/http"

	_ "api/docs"
	"api/internal/domain/session"
	"api/internal/transport/http/helper"

	"github.com/go-playground/validator/v10"
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
// @Failure      400      {object}  helper.ResponseError
// @Failure      401      {object}  helper.ResponseError
// @Failure      500      {object}  helper.ResponseError

type httpServer struct {
	e           *echo.Echo
	sessionRepo session.Repository
}

func New(sessionRepo session.Repository) *httpServer {
	return &httpServer{
		e:           echo.New(),
		sessionRepo: sessionRepo,
	}
}

func (h *httpServer) Run(addr string) error {
	return h.e.Start(addr)
}

func (h *httpServer) Shutdown(c context.Context) error {
	return h.e.Shutdown(c)
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func (h *httpServer) SetupRouter(handlers ...helper.Handler) {
	h.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	h.e.Validator = &CustomValidator{validator: validator.New()}

	h.e.GET("/swagger/*", echoSwagger.WrapHandler)

	baseUrl := h.e.Group("/api/v1")

	for _, handler := range handlers {
		handler.Register(baseUrl, h)
	}
}
